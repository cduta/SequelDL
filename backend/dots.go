package backend

import (
  "database/sql"
)

type Dots struct {
  *Rows
}

type Dot struct {
  Object
  Position
  Color 
  Id int64
}

func InsertDot(handle *Handle, pos Position, color Color) (int64, error) {
  var (
    err          error
    result       sql.Result
    lastInsertId int64
  ) 

  result, err = handle.exec(`
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

  result, err = handle.exec(`
INSERT OR ROLLBACK INTO dots(object_id, x, y) VALUES (?, ?, ?);
INSERT OR ROLLBACK INTO colors(object_id, r, g, b, a) VALUES (?, ?, ?, ?, ?);
COMMIT;
`, lastInsertId, pos.X, pos.Y, lastInsertId, color.R, color.G, color.B, color.A)
  
  return lastInsertId, err
}

func (handle *Handle) QueryDots() (*Dots, error) {
  var (
    err   error 
    rows *Rows
  )

  rows, err = handle.queryRows(`
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
    object   Object   = Object{}
    position Position = Position{}
    color    Color    = Color{}
    dotId    int64
  )

  if !dots.Rows.next() {
    return nil, err
  }

  err = dots.Rows.rows.Scan(&object.Id, &dotId, &position.X, &position.Y, &color.R, &color.G, &color.B, &color.A)

  return &Dot{ Object: object, Position: position, Color: color, Id: dotId }, err
}
