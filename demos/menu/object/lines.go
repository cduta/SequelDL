package object

import (
  "database/sql"

  "../../../backend"
)

type Lines struct {
  *backend.Rows
}

type Line struct {
  backend.Object
  backend.Color 
  Id     int64
  Here   backend.Position
  There  backend.Position
}

func InsertLine(handle *backend.Handle, here backend.Position, there backend.Position, color backend.Color) (int64, error) {
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

  _, err = handle.Exec(`
INSERT OR ROLLBACK INTO lines(object_id, here_x, here_y, there_x, there_y) VALUES (?, ?, ?, ?, ?);
INSERT OR ROLLBACK INTO colors(object_id, r, g, b, a) VALUES (?, ?, ?, ?, ?);
COMMIT;
`, lastInsertId, here.X, here.Y, there.X, there.Y, lastInsertId, color.R, color.G, color.B, color.A)

  return lastInsertId, err
}

func UpdateLineThere(handle *backend.Handle, lineId int64, there backend.Position) error {
  var err error 

  _, err = handle.Exec(`
UPDATE lines   
SET    there_x = ?, there_y = ? 
WHERE  object_id = ?
`, there.X, there.Y, lineId)

  return err
}

func QueryLines(handle *backend.Handle) (*Lines, error) {
  var (
    err   error 
    rows *backend.Rows
  )

  rows, err = handle.QueryRows(`
SELECT l.id, l.object_id, l.here_x, l.here_y, l.there_x, l.there_y, c.r, c.g, c.b, c.a
FROM   lines AS l, colors AS c 
WHERE  l.object_id = c.object_id;
`)

  if rows == nil || err != nil {
    return nil, err 
  }

  return &Lines{ Rows: rows }, err
}

func (lines Lines) Close() {
  lines.Rows.Close()
}

func (lines Lines) Next() (*Line, error) {
  var (
    err      error
    object   backend.Object   = backend.Object{}
    color    backend.Color    = backend.Color{}
    lineId   int64
    here     backend.Position = backend.Position{}
    there    backend.Position = backend.Position{}
  )

  if !lines.Rows.Next() {
    return nil, err
  }

  err = lines.Rows.Scan(
    &lineId , &object.Id, 
    &here.X , &here.Y, 
    &there.X, &there.Y, 
    &color.R, &color.G, &color.B, &color.A)
  if err != nil {
    return nil, err 
  }

  return &Line{ Object: object, Color: color, Id: lineId, Here: here, There: there }, err
}

