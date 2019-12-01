package backend

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

