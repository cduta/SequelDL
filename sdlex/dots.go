package sdlex

import (
  "../backend"
)

func (sdlWrap Wrap) RenderDot(dot *backend.Dot) {
  sdlWrap.renderer.SetDrawColor(dot.Color.R, dot.Color.G, dot.Color.B, dot.Color.A)
  sdlWrap.renderer.DrawPoint(dot.Position.X, dot.Position.Y)  
}

func (sdlWrap Wrap) renderDots() error {
  var (
    err   error 
    dots *backend.Dots
    dot  *backend.Dot
  )

  dots, err = sdlWrap.handle.QueryDots()

  if err != nil {
    return err
  }
  defer dots.Close()

  for dot, err = dots.Next(); err == nil && dot != nil; dot, err = dots.Next() {
    sdlWrap.RenderDot(dot)
  }

  return err
}