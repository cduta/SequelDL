package backend

import (
  "database/sql"
)

type Lines struct {
  *Rows
}

type Line struct {
  Object
  Color 
  Id     int64
  Here   Position
  There  Position
}

func InsertLine(handle *Handle, here Position, there Position, color Color) (int64, error) {
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

func UpdateLineThere(handle *Handle, lineId int64, there Position) error {
  var err error 

  _, err = handle.Exec(`
UPDATE lines   
SET    there_x = ?, there_y = ? 
WHERE  object_id = ?
`, there.X, there.Y, lineId)

  return err
}

func (handle *Handle) QueryLines() (*Lines, error) {
  var (
    err   error 
    rows *Rows
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
    object   Object   = Object{}
    color    Color    = Color{}
    lineId   int64
    here     Position = Position{}
    there    Position = Position{}
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

