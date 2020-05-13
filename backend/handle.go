package backend

import (
  "os"
  "fmt"
  "io/ioutil"
  "database/sql"

  _ "github.com/mattn/go-sqlite3"
)

type Save func(handle *Handle) error
type Load func(handle *Handle) error 

type Handle struct {
  save      Save
  load      Load
  dbhandle *sql.DB
  locked    bool 
}

type Rows struct {
  rows   *sql.Rows
  handle *Handle 
}

type Row struct {
  row    *sql.Row 
}

type Position struct {
  X, Y int32
}

func (position Position) Equals(other Position) bool {
  return position.X == other.X && position.Y == other.Y
}

type Object struct { 
  Id int64
}

func MakeHandle(save Save, load Load, saveFilePath string) (*Handle, error) {
  var (
    err       error
    dbhandle *sql.DB
    loaded    bool
    handle   *Handle
  )
  
  dbhandle, err = sql.Open("sqlite3", "file:backend.db?cache=shared&mode=memory&_foreign_keys=true")
  if err != nil {
    return nil, err
  }

  handle = &Handle{
    save    : save,
    load    : load,
    dbhandle: dbhandle,
    locked  : false}

  err = handle.runSQLFile("backend/initialize.sql")
  if err != nil {
    dbhandle.Close()
    return nil, err
  }

  if saveFilePath == "" {
    err = handle.runSQLFile("backend/options.sql")
    if err != nil {
      dbhandle.Close()
      return nil, err
    }

    err = handle.runSQLFile("backend/ressources.sql")
    if err != nil {
      dbhandle.Close()
      return nil, err
    }
  } else {
    loaded, err = handle.Load(saveFilePath)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Failed to load from file with an error (%s): %s\n", saveFilePath, err)
      dbhandle.Close()
      return nil, err
    } else if !loaded {
      fmt.Fprintf(os.Stderr, "Could not load from file: %s\n", saveFilePath)
    }
  }

  return handle, err
}

func (handle *Handle) Close() {
  handle.dbhandle.Close()
}

func (handle *Handle) runSQLFile(relativeFilePath string) error {
  var (
    err                error
    initQuery        []byte
    absoluteFilePath   string 
  )

  absoluteFilePath, err = ToAbsolutePath(relativeFilePath)
  if err != nil {
    return err
  }

  _, err = os.Stat(absoluteFilePath)
  if err != nil {
    return fmt.Errorf("Could not run SQL file at %s\n", absoluteFilePath) 
  }

  initQuery, err = ioutil.ReadFile(absoluteFilePath)
  if err != nil {
    return err
  }

  _, err = handle.dbhandle.Exec(string(initQuery))
  if err != nil {
    return err
  }

  return err
}

func (handle *Handle) queryRow(query string, args ...interface{}) (*Row, error) {
  var (
    err  error 
    row *sql.Row
  )

  if !handle.isLocked() {
    row, err = handle.dbhandle.QueryRow(query, args...), err 

    if err == sql.ErrNoRows {
      return nil, nil
    }
  }

  return &Row{row: row}, err
}

func (handle *Handle) query(query string, args ...interface{}) (*sql.Rows, error) {
  var (
    err   error 
    rows *sql.Rows
  )

  if !handle.isLocked() {
    handle.lock()
    return handle.dbhandle.Query(query, args...)
  }

  return rows, err
}

func (handle *Handle) Exec(query string, args ...interface{}) (sql.Result, error) {
  var (
    err    error 
    result sql.Result
  )

  if !handle.isLocked() {
    return handle.dbhandle.Exec(query, args...)
  }

  return result, err
}

func (handle *Handle) Save(path string) (bool, error) {
  var (
    err          error
    saveHandle  *sql.DB
    schemaRows  *sql.Rows
    statement    string   
    schema       string 
  )

  if handle.isLocked() {
    return false, err
  }

  schemaRows, err = handle.query(`
    SELECT sm.sql 
    FROM   sqlite_master AS sm
    WHERE  sm.name NOT LIKE 'sqlite_%';
  `)
  if err != nil {
    return false, err 
  }
  defer schemaRows.Close()

  for schemaRows.Next() {
    err = schemaRows.Scan(&statement)
    if err != nil {
      return false, err
    }

    schema = schema + statement + ";\n"
  }

  schemaRows.Close()

  _, err = os.Stat("save.db")
  if err == nil {
    os.Remove("save.db")
  }

  saveHandle, err = sql.Open("sqlite3", "file:"+path+"?cache=shared&_foreign_keys=true")
  if err != nil {
    return false, err
  }
  defer saveHandle.Close()

  _, err = saveHandle.Exec(schema)
  if err != nil {
    return false, err
  }

  saveHandle.Close()

  _, err = handle.Exec(`
ATTACH DATABASE ? AS save;
BEGIN IMMEDIATE;
`, "file:"+path+"?cache=shared&_foreign_keys=true")
  if err != nil {
    return false, err
  }

  err = handle.save(handle)
  if err != nil {
    return false, err
  }

  _, err = handle.Exec(`
COMMIT;
DETACH DATABASE save;
`)
  if err != nil {
    return false, err
  }

  return true, err
}

func (handle *Handle) Load(path string) (bool, error) {
  var err error

  _, err = handle.Exec(`
ATTACH DATABASE ? AS save;
BEGIN IMMEDIATE;
`, "file:"+path+"?cache=shared&_foreign_keys=true")
  if err != nil {
    return false, err
  }

  err = handle.load(handle)
  if err != nil {
    return false, err
  }

  _, err = handle.Exec(`
COMMIT;
DETACH DATABASE save;
`)
  if err != nil {
    return false, err
  }

  return true, err
}

func (handle *Handle) lock() {
  handle.locked = true
}

func (handle *Handle) unlock() {
  handle.locked = false
}

func (handle *Handle) isLocked() bool {
  return handle.locked
}

func (handle *Handle) queryRows(query string, args ...interface{}) (*Rows, error) {
  var ( 
    err   error
    rows *sql.Rows
  ) 

  rows, err = handle.query(query, args...)

  if err != nil {
    return nil, err 
  }

  return &Rows{ rows: rows, handle: handle }, err
}

func (rows *Rows) next() bool {
  var hasNext bool = false 

  if rows.rows != nil {
    hasNext = rows.rows.Next()
  }
  return hasNext
}

func (rows *Rows) Close() {
  rows.handle.unlock()
  if rows.rows != nil {
    rows.rows.Close()
  }
}

func (rows *Rows) Scan(args ...interface{}) error {
  return rows.rows.Scan(args...)
}

func (row *Row) Scan(args ...interface{}) error {
  return row.row.Scan(args...)
}