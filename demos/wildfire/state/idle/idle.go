package idle

import (
  "../../../../backend"
  "../../../../sdlex"
  . "../../../../event/state"
  "../../wrap"

  "github.com/veandco/go-sdl2/sdl"

  "fmt"
  "os"
)

type Idle struct { 
  particles     *wrap.Particles
  backendHandle *backend.Handle 
}

func MakeIdle(particles *wrap.Particles, backendHandle *backend.Handle) (Idle, error) { 
  var err error

  err = particles.LoadParticles(backendHandle)

  return Idle{ 
    particles    : particles,
    backendHandle: backendHandle }, err
}

func (idle Idle) Destroy() {}
func (idle Idle) PreEvent()    State { return idle }

func (idle Idle) OnTick()      State {
  //idle.particles.Animate()  
  return idle 
}

func (idle Idle) TickDelayed() bool  { return false }
func (idle Idle) OnTickDelay() State { return idle }

func (idle Idle) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(idle) }
func (idle Idle) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var (
    err   error
    state State = idle
    saved bool
  )

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      if event.Keysym.Mod & sdl.KMOD_CTRL > 0 {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(idle)
          case sdl.K_s:
            saved, err = idle.backendHandle.Save("save.db")
            if err != nil {
              fmt.Fprintf(os.Stderr, "Failed to load file with an error: %s\n", err)
            } else if !saved {
              fmt.Fprintf(os.Stderr, "Could not save to file. Try again.")              
            }
        }
      } else {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(idle)
        }
      }
  } 

  return state
}

func (idle Idle) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State { return idle }
func (idle Idle) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State { return idle }
func (idle Idle) PostEvent() State { return idle }