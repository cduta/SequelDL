package sdlex

import (
  "../backend"
)

func (sdlWrap Wrap) RenderLine(line *backend.Line) {
  var part backend.Position

  sdlWrap.renderer.SetDrawColor(line.Color.R, line.Color.G, line.Color.B, line.Color.A)

  for _, part = range line.Parts {
    sdlWrap.renderer.DrawPoint(part.X, part.Y)
  }
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