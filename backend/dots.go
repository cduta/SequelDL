package backend

type Dots struct {
  *Objects
}

type Dot struct {
  Object
  Position
  Color 
  Id int64
}

func (handle Handle) QueryDots() (*Dots, error) {
  var (
    err      error 
    objects *Objects
  )

  objects, err = handle.queryObjects(`
SELECT o.id, d.id, d.x, d.y, c.r, c.g, c.b, c.a
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