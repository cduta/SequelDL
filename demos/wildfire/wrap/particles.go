package wrap

import (
  "../object"
  "../../../backend"
  "../../../sdlex"

  "github.com/veandco/go-sdl2/sdl"

  "fmt"
  "math/rand"
)

type Particle struct {
  *object.Particle
  Color backend.Color
}

type Particles struct {
  byPosition  map[backend.Position]*Particle
  randomizer *rand.Rand
}

func (particles *Particles) shuffleColorRange(current *backend.Color, from backend.Color, to backend.Color) {
  current.R = from.R + uint8(particles.randomizer.Uint32() % (uint32(to.R-from.R+1)))
  current.G = from.G + uint8(particles.randomizer.Uint32() % (uint32(to.G-from.G+1)))
  current.B = from.B + uint8(particles.randomizer.Uint32() % (uint32(to.B-from.B+1)))
  current.A = from.A + uint8(particles.randomizer.Uint32() % (uint32(to.A-from.A+1)))
}

func (particles *Particles) LoadParticles(handle *backend.Handle) error {
  var (
    err              error 
    objectParticles *object.Particles
    objectParticle  *object.Particle
    particle        *Particle
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
    particle = &Particle{ 
      Particle: objectParticle, 
      Color   : backend.Color{}}
    particles.shuffleColorRange(&particle.Color, objectParticle.From, objectParticle.To)
    particlesByPos[objectParticle.Position] = particle
  }

  particles.byPosition = particlesByPos

  return err 
}

func (particles *Particles) UpdateParticles(handle *backend.Handle) error {
  var (
    err               error 
    oldParticles     *object.Particles
    changedParticles *object.Particles
    objectParticle   *object.Particle
    updatedParticle  *Particle
  )

  oldParticles, err = object.QueryOldParticles(handle)
  if oldParticles == nil || err != nil {
    return err
  }
  defer oldParticles.Close() 

  for objectParticle, err = oldParticles.Next(); err == nil && objectParticle != nil; objectParticle, err = oldParticles.Next() {
    if err != nil {
      return err;
    }
    delete(particles.byPosition, objectParticle.Position)
  }

  oldParticles.Close() 

  changedParticles, err = object.QueryChangedParticles(handle)
  if changedParticles == nil || err != nil {
    return err
  }
  defer changedParticles.Close()

  for objectParticle, err = changedParticles.Next(); err == nil && objectParticle != nil; objectParticle, err = changedParticles.Next() {
    if err != nil {
      return err;
    }
    updatedParticle = &Particle{ 
      Particle: objectParticle, 
      Color   : backend.Color{}}
    particles.shuffleColorRange(&updatedParticle.Color, updatedParticle.Particle.From, updatedParticle.Particle.To) 
    particles.byPosition[updatedParticle.Particle.Position] = updatedParticle
  }

  changedParticles.Close()

  return object.ClearOldStates(handle) 
} 

func (particles *Particles) Animate() {
  var (
    redraw    bool
    particle *Particle 
  )

  for _, particle = range particles.byPosition {
    redraw = particle.AdvanceRedrawDelay()
    if redraw {
       particles.shuffleColorRange(&particle.Color, particle.Particle.From, particle.Particle.To) 
    }
  }
}

func (wildfireWrap *WildfireWrap) RenderParticlesByColor(sdlWrap *sdlex.SdlWrap, color backend.Color, points []sdl.Point) error {
  var err error

  err = sdlWrap.Renderer().SetDrawColor(color.R, color.G, color.B, color.A)
  if err != nil {
    return err 
  }

  return sdlWrap.Renderer().DrawPoints(points)
}

func (wildfireWrap *WildfireWrap) RenderParticles(sdlWrap *sdlex.SdlWrap) error {
  var (
    err                error 
    position           backend.Position
    particle          *Particle 
    particlesByColor   map[backend.Color][]sdl.Point = make(map[backend.Color][]sdl.Point)  
    color              backend.Color 
    points           []sdl.Point  
  )

  for position, particle = range wildfireWrap.Particles().byPosition {
    particlesByColor[particle.Color] = append(particlesByColor[particle.Color], 
                                              sdl.Point{X: position.X, 
                                                        Y: position.Y})
  }


  for color, points = range particlesByColor {
    err = wildfireWrap.RenderParticlesByColor(sdlWrap, color, points)
    if err != nil {
      fmt.Errorf("Could not render particle at %v with color %v.", points, particle.Color)
      return err
    }
  }

  return err
}