package wrap

import (
  "../object"
  "../../../backend"
  "../../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

func (wildfireWrap *WildfireWrap) RenderParticle(sdlWrap *sdlex.SdlWrap, color backend.Color, points []sdl.Point) error {
  var (
    err       error
    renderer *sdl.Renderer = sdlWrap.Renderer()
  ) 

  err = renderer.SetDrawColor(color.R, color.G, color.B, color.A)
  if err != nil {
    return err 
  }

  err = renderer.DrawPoints(points)

  return err
}

func (wildfireWrap *WildfireWrap) RenderParticles(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error {
  var (
    err          error 
    particles   *object.Particles
    particle    *object.Particle
    particleMap  map[backend.Color][]sdl.Point = make(map[backend.Color][]sdl.Point)
    color        backend.Color 
    points       []sdl.Point
  )

  particles, err = object.QueryParticles(handle)
  if particles == nil || err != nil {
    return err
  }
  defer particles.Close()

  for particle, err = particles.Next(); err == nil && particle != nil; particle, err = particles.Next() {
    if err != nil {
      return err
    }
    particleMap[backend.Color{particle.R, particle.G, particle.B, particle.A}] = append(
      particleMap[backend.Color{R:particle.Color.R, G:particle.Color.G, B:particle.Color.B, A:particle.Color.A}], 
      sdl.Point{X:particle.Position.X, Y:particle.Position.Y})
  }

  for color, points = range particleMap {
    err = wildfireWrap.RenderParticle(sdlWrap, color, points)
    if err != nil {
      return err
    }
  }

  return err
}