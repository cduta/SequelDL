package backend

import (
  "os"
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

type Position struct {
  X, Y int32
}

type Object struct { 
  Id int64
}

type Color struct {
  R,G,B,A uint8 
}

type Dots struct {
  *Objects
}

type Dot struct {
  Object
  Position
  Color 
  Id int64
}

type Lines struct {
  *Objects
}

type Line struct {
  Object
  Color 
  Id    int64
  Here  Position
  There Position
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

func (handle Handle) Query(query string, args ...interface{}) (*sql.Rows, error) {
  return handle.dbhandle.Query(query, args...)
}

func (handle Handle) Exec(query string, args ...interface{}) (sql.Result, error) {
  return handle.dbhandle.Exec(query, args...)
}

func (handle Handle) Save(path string) error {
  var (
    err          error
    saveHandle  *sql.DB
    schemaRows  *sql.Rows
    statement    string   
    schema       string 
  )

  schemaRows, err = handle.dbhandle.Query(`
    SELECT sm.sql 
    FROM   sqlite_master AS sm
    WHERE  sm.name NOT LIKE 'sqlite_%';
  `)
  if err != nil {
    return err 
  }
  defer schemaRows.Close()

  for schemaRows.Next() {
    err = schemaRows.Scan(&statement)
    if err != nil {
      return err
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
    return err
  }
  defer saveHandle.Close()

  _, err = saveHandle.Exec(schema)
  if err != nil {
    return err
  }

  saveHandle.Close()

  _, err = handle.dbhandle.Exec(`
ATTACH DATABASE ? AS save;

BEGIN IMMEDIATE;
INSERT OR ROLLBACK INTO save.objects          SELECT * FROM main.objects;
INSERT OR ROLLBACK INTO save.dots             SELECT * FROM main.dots;
INSERT OR ROLLBACK INTO save.lines            SELECT * FROM main.lines;
INSERT OR ROLLBACK INTO save.rectangles       SELECT * FROM main.rectangles;
INSERT OR ROLLBACK INTO save.triangles        SELECT * FROM main.triangles;
INSERT OR ROLLBACK INTO save.polygons         SELECT * FROM main.polygons;
INSERT OR ROLLBACK INTO save.polygon_vertices SELECT * FROM main.polygon_vertices;
INSERT OR ROLLBACK INTO save.colors           SELECT * FROM main.colors;
COMMIT;

DETACH DATABASE save;
`, "file:"+path+"?cache=shared&_foreign_keys=true")
  if err != nil {
    return err
  }

  return err
}

func (handle Handle) queryObjects(query string, args ...interface{}) (*Objects, error) {
  var ( 
    err   error
    rows *sql.Rows
  ) 

  rows, err = handle.dbhandle.Query(query, args...)

  if err != nil {
    return nil, err 
  }

  return &Objects{ rows: rows }, err
}

func (objects *Objects) Close() {
  objects.rows.Close()
}

func (handle Handle) QueryDots() (*Dots, error) {
  var (
    err      error 
    objects *Objects
  )

  objects, err = handle.queryObjects(`
SELECT o.id, d.x, d.y, c.r, c.g, c.b, c.a
FROM   objects AS o, dots AS d, colors AS c  
WHERE  o.id = d.object_id
AND    o.id = c.object_id;`)

  if err != nil {
    return nil, err 
  }

  return &Dots{ Objects: objects }, err
}

func (dots Dots) Close() {
  dots.Objects.Close()
}

func (dots Dots) Next() (*Dot, error) {
  var (
    err      error
    object   Object   = Object{}
    position Position = Position{}
    color    Color    = Color{}
    dotId    int64
  )

  if !dots.Objects.rows.Next() {
    return nil, err
  }

  err = dots.Objects.rows.Scan(&object.Id, &dotId, &position.X, &position.Y, &color.R, &color.G, &color.B, &color.A)

  return &Dot{ Object: object, Position: position, Color: color, Id: dotId }, err
}

func (handle Handle) QueryLines() (*Lines, error) {
  var (
    err      error 
    objects *Objects
  )

  objects, err = handle.queryObjects(`
SELECT o.id, l.id, l.here_x, l.here_y, l.there_x, l.there_y, c.r, c.g, c.b, c.a
FROM   objects AS o, lines AS l, colors AS c  
WHERE  o.id = l.object_id
AND    o.id = c.object_id;`)

  if err != nil {
    return nil, err 
  }

  return &Lines{ Objects: objects }, err
}

func (lines Lines) Close() {
  lines.Objects.Close()
}

func (lines Lines) Next() (*Line, error) {
  var (
    err    error
    object Object   = Object{}
    color  Color    = Color{}
    lineId int64
    here   Position = Position{}
    there  Position = Position{}
  )

  if !lines.Objects.rows.Next() {
    return nil, err
  }

  err = lines.Objects.rows.Scan(&object.Id, &lineId, &here.X, &here.Y, &there.X, &there.Y, &color.R, &color.G, &color.B, &color.A)

  return &Line{ Object: object, Color: color, Id: lineId, Here: here, There: there }, err
}
