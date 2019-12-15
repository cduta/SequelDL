package backend

import (
  "errors"
  "fmt"
  "io/ioutil"
  "path/filepath"
  "database/sql"
)

type Lines struct {
  *Objects
}

type Line struct {
  Object
  Color 
  Id      int64
  Parts []Position
}

func InsertLine(handle *Handle, here Position, there Position, color Color) (int64, error) {
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

  _, err = handle.exec(`
INSERT OR ROLLBACK INTO lines(object_id, here_x, here_y, there_x, there_y) VALUES (?, ?, ?, ?, ?);
INSERT OR ROLLBACK INTO colors(object_id, r, g, b, a) VALUES (?, ?, ?, ?, ?);
COMMIT;
`, lastInsertId, here.X, here.Y, there.X, there.Y, lastInsertId, color.R, color.G, color.B, color.A)

  return lastInsertId, err
}

func UpdateLineThere(handle *Handle, lineId int64, there Position) error {
  var err error 

  _, err = handle.exec(`
UPDATE lines   
SET    there_x = ?, there_y = ? 
WHERE  object_id = ?
`, there.X, there.Y, lineId)

  return err
}

func (handle Handle) QueryLines() (*Lines, error) {
  var (
    err         error 
    bresenham []byte
    objects    *Objects
  )

  bresenham, err = ioutil.ReadFile(filepath.Join("backend", "sql", "bresenham.sql"))
  if err != nil {
    return nil, err
  }

  objects, err = handle.queryObjects(string(bresenham))

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
    err      error
    object   Object   = Object{}
    color    Color    = Color{}
    lineId   int64
    curr     Position = Position{}
    there    Position = Position{}
    parts  []Position
  )

  if !lines.Objects.rows.Next() {
    return nil, err
  }

  for lines.Objects.rows.Next() {
    err = lines.Objects.rows.Scan(
      &lineId , &object.Id, 
      &curr.X , &curr.Y, 
      &there.X, &there.Y, 
      &color.R, &color.G, &color.B, &color.A)
    if err != nil {
      return nil, err 
    }

    parts = append(parts, curr)

    if curr.Equals(there) {
      break
    }
  }

  if !curr.Equals(there) {
    return nil, errors.New(fmt.Sprintf("Queried line (Id: %d, Object Id: %d) incomplete", lineId, object.Id))
  }

  return &Line{ Object: object, Color: color, Id: lineId, Parts: parts }, err
}

