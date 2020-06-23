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

type Burn struct { 
  particles     *wrap.Particles
  backendHandle *backend.Handle 
}

func MakeIdle(particles *wrap.Particles, backendHandle *backend.Handle) (Burn, error) { 
  var err error

  return Burn{ 
    particles    : particles,
    backendHandle: backendHandle }, err
}

func (burn Burn) Destroy() {}
func (burn Burn) PreEvent()    State { return burn }

func (burn Burn) OnTick()      State {
  //object.BurnPaper(...)  
  return burn 
}

func (burn Burn) TickDelayed() bool  { return false }
func (burn Burn) OnTickDelay() State { return burn }

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
          burn.particles.UpdateParticles(burn.backendHandle)
          if err != nil {
            fmt.Fprintf(os.Stderr, "Error when trying to ignite paper: %s\n", err)
            return burn  
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

