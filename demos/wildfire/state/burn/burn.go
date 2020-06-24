package burn

import (
  "../../../../backend"
  "../../../../sdlex"
  . "../../../../event/state"
  "../../wrap"
  "../../object"

  "github.com/veandco/go-sdl2/sdl"

  "fmt"
  "os"
)

const TICKS_UNTIL_BURN = 2

type Burn struct { 
  particles     *wrap.Particles
  backendHandle *backend.Handle 
  burnTicks      int64
}

func MakeIdle(particles *wrap.Particles, backendHandle *backend.Handle) (Burn, error) { 
  var err error

  return Burn{ 
    particles    : particles,
    backendHandle: backendHandle,
    burnTicks    : TICKS_UNTIL_BURN }, err
}

func (burn Burn) Destroy() {}
func (burn Burn) PreEvent()    State { return burn }

func (burn Burn) OnTick()      State {
  var (
    err   error 
    burns int64 
  )

  burns, err = object.BurnPaper(burn.backendHandle)  
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error when trying to burn paper: %s\n", err)
  }

  if err == nil && burns > 0 {
    err = burn.particles.UpdateParticles(burn.backendHandle)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error when trying to update particles: %s\n", err)
    }
  }

  burn.burnTicks = TICKS_UNTIL_BURN

  return burn 
}

func (burn Burn) TickDelayed() bool  { 
  return burn.burnTicks >= 0
}

func (burn Burn) OnTickDelay() State { 
  burn.burnTicks--
  return burn 
}

func (burn Burn) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(burn) }
func (burn Burn)    OnKeyboardEvent(event *sdl.KeyboardEvent)    State { return burn }
func (burn Burn) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State { return burn }
func (burn Burn) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State { 
  var (
    err       error
    state     State = burn
  )

  switch event.State {
    case sdlex.BUTTON_PRESSED: 
      switch event.Button {
        case sdl.BUTTON_LEFT  : 
          err = object.IgnitePaper(burn.backendHandle, backend.Position{X: event.X, Y: event.Y})
          if err != nil {
            fmt.Fprintf(os.Stderr, "Error when trying to ignite paper: %s\n", err)
            return burn  
          }
          err = burn.particles.UpdateParticles(burn.backendHandle)
          if err != nil {
            fmt.Fprintf(os.Stderr, "Error when trying to update particles: %s\n", err)
          }
        case sdl.BUTTON_RIGHT : 
        default               : 
      }
    case sdlex.BUTTON_RELEASED: 
      switch event.Button {
        case sdl.BUTTON_LEFT  :  
        case sdl.BUTTON_RIGHT : 
        default               :
      }
  }

  return state 
}
func (burn Burn) PostEvent() State { return burn }

