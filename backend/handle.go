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
  Id           int
  ObjectTypeId int
  Name         string
  X, Y         int32 
}

func NewHandle() (*Handle, error) {
  var (
    err      error
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
    err       error
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

func (handle Handle) AddObject(name string, x int32, y int32) error {
  var err error 

  _, err = handle.dbhandle.Exec(`
INSERT INTO objects(objecttype_id,x,y) 
SELECT ot.id, ?, ?
FROM   objecttype AS ot
WHERE  ot.name = ?;`, x, y, name)
  if err != nil {
    return err
  }

  return err
}

func (handle Handle) QueryObjects(name string) (*Objects, error) {
  var ( 
    err  error
    rows *sql.Rows
  ) 

  rows, err = handle.dbhandle.Query(`
SELECT o.id, ot.id, ot.name, o.x, o.y
FROM   objects AS o
JOIN   objecttype AS ot 
ON     o.objecttype_id = ot.id 
WHERE  ot.name = ?;`, name)

  if err != nil {
    return nil, err 
  }

  return &Objects{ rows: rows }, err
}

func (objects Objects) Close() {
  objects.rows.Close()
}

func (objects Objects) NextObject() (*Object, error) {
  var (
    err error
    object *Object = &Object{}
  )

  if !objects.rows.Next() {
    return nil, err
  }

  err = objects.rows.Scan(&object.Id, &object.ObjectTypeId, &object.Name, &object.X, &object.Y)

  return object, err
}

