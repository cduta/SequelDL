package wrap

import (
  "../object"
  "../../../backend"
  "../../../sdlex"

  "github.com/veandco/go-sdl2/gfx"
)

func (menuWrap *MenuWrap) RenderDot(sdlWrap *sdlex.SdlWrap, dot *object.Dot) {
  gfx.PixelRGBA(
    sdlWrap.Renderer(), 
    dot.Position.X, dot.Position.Y,
    dot.Color.R, dot.Color.G, dot.Color.B, dot.Color.A)
}

func (menuWrap *MenuWrap) RenderDots(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error {
  var (
    err   error 
    dots *object.Dots
    dot  *object.Dot
  )

  dots, err = object.QueryDots(handle)
  if dots == nil || err != nil {
    return err
  }
  defer dots.Close()

  for dot, err = dots.Next(); err == nil && dot != nil; dot, err = dots.Next() {
    menuWrap.RenderDot(sdlWrap, dot)
  }

  return err
}