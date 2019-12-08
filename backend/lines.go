package backend

import (
  "errors"
  "fmt"
  "io/ioutil"
  "path/filepath"
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

