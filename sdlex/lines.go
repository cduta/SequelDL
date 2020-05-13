package sdlex

import (
  "../backend"

  "github.com/veandco/go-sdl2/gfx"
)

func (sdlWrap SdlWrap) RenderLine(line *backend.Line) {
  gfx.PolygonRGBA(
    sdlWrap.renderer, 
    []int16{int16(line.Here.X), int16(line.There.X), int16(line.Here.X)},  
    []int16{int16(line.Here.Y), int16(line.There.Y), int16(line.Here.Y)}, 
    line.Color.R, line.Color.G, line.Color.B, line.Color.A)
}

func (sdlWrap SdlWrap) RenderLines() error {
  var (
    err    error 
    lines *backend.Lines
    line  *backend.Line
  )

  lines, err = sdlWrap.handle.QueryLines()
  if lines == nil || err != nil {
    return err
  }
  defer lines.Close()

  for line, err = lines.Next(); err == nil && line != nil; line, err = lines.Next() {
    sdlWrap.RenderLine(line)
  }

  return err
}