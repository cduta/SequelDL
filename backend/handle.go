package backend

import (
  "io/ioutil"
  "path/filepath"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type Handle struct {
  dbhandle *sql.DB
}

type Objects struct {
  rows *sql.Rows
}

type Object struct {
  Id   int
  X, Y int32 
}

type Dots struct {
	*Objects
}

type Dot struct {
	Object
}

func NewHandle() (*Handle, error) {
  var (
    err       error
    dbhandle *sql.DB
  )
  
  dbhandle, err = sql.Open("sqlite3", "file:backend.db?cache=shared&mode=memory&_foreign_keys=true")
  if err != nil {
    return nil, err
  }

  err = initializeBackend(dbhandle)
  if err != nil {
    return nil, err
  }

  return &Handle{
    dbhandle: dbhandle}, err
}

func (handle *Handle) Close() {
  handle.dbhandle.Close()
}

func initializeBackend(dbhandle *sql.DB) error {
  var (
    err         error
    initQuery []byte
  )

  initQuery, err = ioutil.ReadFile(filepath.Join("backend", "sql", "initializeBackend.sql"))
  if err != nil {
    return err
  }

  _, err = dbhandle.Exec(string(initQuery))
  if err != nil {
    return err
  }

  return err
}

func (handle Handle) AddDot(x int32, y int32) error {
  var err error 

  _, err = handle.dbhandle.Exec(`
INSERT INTO object(x,y)    VALUES (?, ?);
INSERT INTO dot(object_id) VALUES (last_insert_rowid());
`, x, y)
  
  return err
}

func (handle Handle) queryObjects(query string) (*Objects, error) {
  var ( 
    err   error
    rows *sql.Rows
  ) 

  rows, err = handle.dbhandle.Query(query)

  if err != nil {
    return nil, err 
  }

  return &Objects{ rows: rows }, err
}

func (handle Handle) QueryDots() (*Dots, error) {
	var (
		err      error 
		objects *Objects
	)

	objects, err = handle.queryObjects(`
SELECT o.id, o.x, o.y
FROM   object AS o JOIN dot AS d 
ON     o.id = d.object_id;`)

	if err != nil {
		return nil, err 
	}

	return &Dots{ Objects: objects }, err
}

func (objects *Objects) Close() {
	objects.rows.Close()
}

func (dots Dots) Close() {
  dots.Objects.Close()
}

func (dots Dots) Next() (*Dot, error) {
  var (
    err    error
    object Object = Object{}
  )

  if !dots.Objects.rows.Next() {
    return nil, err
  }

  err = dots.Objects.rows.Scan(&object.Id, &object.X, &object.Y)

  return &Dot{ Object: object }, err
}

