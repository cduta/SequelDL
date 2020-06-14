package wrap

import (
  "../object"
  "../../../backend"
  "../../../sdlex"

  "github.com/veandco/go-sdl2/gfx"

  "fmt"
  "math/rand"
)

type Particle struct {
  *object.Particle
  Color backend.Color
}

type Particles struct {
  byPosition map[backend.Position]*Particle
}

func (particles *Particles) LoadParticles(handle *backend.Handle) error {
  var (
    err              error 
    objectParticles *object.Particles
    objectParticle  *object.Particle
    particlesByPos   map[backend.Position]*Particle = make(map[backend.Position]*Particle)
  )

  objectParticles, err = object.QueryVisibleParticles(handle)
  if objectParticles == nil || err != nil {
    return err
  }
  defer objectParticles.Close()

  for objectParticle, err = objectParticles.Next(); err == nil && objectParticle != nil; objectParticle, err = objectParticles.Next() {
    if err != nil {
      return err
    }
    particlesByPos[objectParticle.Position] = &Particle{ 
      Particle: objectParticle, 
      Color: backend.Color{
        R: objectParticle.From.R + uint8(rand.Intn(int(objectParticle.To.R-objectParticle.From.R+1))),
        G: objectParticle.From.G + uint8(rand.Intn(int(objectParticle.To.G-objectParticle.From.G+1))),
        B: objectParticle.From.B + uint8(rand.Intn(int(objectParticle.To.B-objectParticle.From.B+1))),
        A: objectParticle.From.A + uint8(rand.Intn(int(objectParticle.To.A-objectParticle.From.A+1)))}}
  }

  particles.byPosition = particlesByPos

  return err 
}

func (wildfireWrap *WildfireWrap) RenderParticle(sdlWrap *sdlex.SdlWrap, position backend.Position, color backend.Color) bool {
  return gfx.PixelRGBA(sdlWrap.Renderer(), position.X, position.Y, color.R, color.G, color.B, color.A)
}

func (wildfireWrap *WildfireWrap) RenderParticles(sdlWrap *sdlex.SdlWrap) error {
  var (
    err      error 
    position backend.Position
    particle *Particle 
  )

  for position, particle = range wildfireWrap.Particles().byPosition {
    if !wildfireWrap.RenderParticle(sdlWrap, position, particle.Color) {
      return fmt.Errorf("Could not render particle at %v with color %v", position, particle.Color)
    }
  }

  return err
}