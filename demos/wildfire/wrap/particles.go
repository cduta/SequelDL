package wrap

import (
  "../object"
  "../../../backend"
  "../../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type Particles struct {
  fire map[backend.Color][]sdl.Point 
}

func (particles *Particles) ReloadParticles(handle *backend.Handle) error {
  var (
    err                error 
    objectParticles   *object.Particles
    objectParticle    *object.Particle
    fireParticleMap    map[backend.Color][]sdl.Point = make(map[backend.Color][]sdl.Point)
  )

  objectParticles, err = object.QueryFireParticles(handle)
  if objectParticles == nil || err != nil {
    return err
  }
  defer objectParticles.Close()

  for objectParticle, err = objectParticles.Next(); err == nil && objectParticle != nil; objectParticle, err = objectParticles.Next() {
    if err != nil {
      return err
    }
    fireParticleMap[backend.Color{objectParticle.R, objectParticle.G, objectParticle.B, objectParticle.A}] = append(
      fireParticleMap[backend.Color{R:objectParticle.Color.R, G:objectParticle.Color.G, B:objectParticle.Color.B, A:objectParticle.Color.A}], 
      sdl.Point{X:objectParticle.Position.X, Y:objectParticle.Position.Y})
  }

  particles.fire = fireParticleMap

  return err
}

func (particles *Particles) HasFireLoaded() bool {
  return particles.fire != nil
}

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

func (wildfireWrap *WildfireWrap) RenderParticles(sdlWrap *sdlex.SdlWrap) error {
  var (
    err      error 
    color    backend.Color 
    points []sdl.Point
  )

  if wildfireWrap.Particles().HasFireLoaded() {
    for color, points = range wildfireWrap.Particles().fire {
      err = wildfireWrap.RenderParticle(sdlWrap, color, points)
      if err != nil {
        return err
      }
    }
  }

  return err
}