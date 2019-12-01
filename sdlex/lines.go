package sdlex

import (
  "../backend"
)

func (sdlWrap Wrap) RenderLine(line *backend.Line) {
  sdlWrap.renderer.SetDrawColor(line.Color.R, line.Color.G, line.Color.B, line.Color.A)
  sdlWrap.renderer.DrawLine(line.Here.X, line.Here.Y, line.There.X, line.There.Y)
}

func (sdlWrap Wrap) renderLines() error {
  var (
    err    error 
    lines *backend.Lines
    line  *backend.Line
  )

  lines, err = sdlWrap.handle.QueryLines()
  if err != nil {
    return err
  }
  defer lines.Close()

  for line, err = lines.Next(); err == nil && line != nil; line, err = lines.Next() {
    sdlWrap.RenderLine(line)
  }

  return err
}