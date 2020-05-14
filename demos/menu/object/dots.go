package object

import (
  "../../../backend"
  "database/sql"
)

type Dots struct {
  *backend.Rows
}

type Dot struct {
  backend.Object
  backend.Position
  backend.Color 
  Id int64
}

func InsertDot(handle *backend.Handle, pos backend.Position, color backend.Color) (int64, error) {
  var (
    err          error
    result       sql.Result
    lastInsertId int64
  ) 

  result, err = handle.Exec(`
BEGIN IMMEDIATE;
INSERT OR ROLLBACK INTO objects DEFAULT VALUES;
`)

  if err != nil {
    return lastInsertId, err 
  }

  lastInsertId, err = result.LastInsertId()

  if err != nil {
    return lastInsertId, err
  }

  result, err = handle.Exec(`
INSERT OR ROLLBACK INTO dots(object_id, x, y) VALUES (?, ?, ?);
INSERT OR ROLLBACK INTO colors(object_id, r, g, b, a) VALUES (?, ?, ?, ?, ?);
COMMIT;
`, lastInsertId, pos.X, pos.Y, lastInsertId, color.R, color.G, color.B, color.A)
  
  return lastInsertId, err
}

func QueryDots(handle *backend.Handle) (*Dots, error) {
  var (
    err   error 
    rows *backend.Rows
  )

  rows, err = handle.QueryRows(`
SELECT o.id, d.id, d.x, d.y, c.r, c.g, c.b, c.a
FROM   objects AS o, dots AS d, colors AS c  
WHERE  o.id = d.object_id
AND    o.id = c.object_id;`)
  if rows == nil || err != nil {
    return nil, err 
  }

  return &Dots{ Rows: rows }, err
}

func (dots Dots) Close() {
  dots.Rows.Close()
}

func (dots Dots) Next() (*Dot, error) {
  var (
    err      error
    object   backend.Object   = backend.Object{}
    position backend.Position = backend.Position{}
    color    backend.Color    = backend.Color{}
    dotId    int64
  )

  if !dots.Rows.Next() {
    return nil, err
  }

  err = dots.Rows.Scan(&object.Id, &dotId, &position.X, &position.Y, &color.R, &color.G, &color.B, &color.A)

  return &Dot{ Object: object, Position: position, Color: color, Id: dotId }, err
}
