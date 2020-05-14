package wrap

import (
  "../object"
  "../../../backend"
  "../../../sdlex"

  "github.com/veandco/go-sdl2/gfx"
)

func (menuWrap *MenuWrap) RenderLine(sdlWrap *sdlex.SdlWrap, line *object.Line) {
  gfx.PolygonRGBA(
    sdlWrap.Renderer(), 
    []int16{int16(line.Here.X), int16(line.There.X), int16(line.Here.X)},  
    []int16{int16(line.Here.Y), int16(line.There.Y), int16(line.Here.Y)}, 
    line.Color.R, line.Color.G, line.Color.B, line.Color.A)
}

func (menuWrap *MenuWrap) RenderLines(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error {
  var (
    err    error 
    lines *object.Lines
    line  *object.Line
  )

  lines, err = object.QueryLines(handle)
  if lines == nil || err != nil {
    return err
  }
  defer lines.Close()

  for line, err = lines.Next(); err == nil && line != nil; line, err = lines.Next() {
    menuWrap.RenderLine(sdlWrap, line)
  }

  return err
}