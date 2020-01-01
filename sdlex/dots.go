package sdlex

import (
  "../backend"

  "github.com/veandco/go-sdl2/gfx"
)

func (sdlWrap Wrap) RenderDot(dot *backend.Dot) {
  gfx.PixelRGBA(
    sdlWrap.renderer, 
    dot.Position.X, dot.Position.Y,
    dot.Color.R, dot.Color.G, dot.Color.B, dot.Color.A)
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