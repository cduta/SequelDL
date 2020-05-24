package idle

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

type Idle struct { 
  stateId        int64
  particles     *wrap.Particles
  backendHandle *backend.Handle 
}

func MakeIdle(stateId int64, particles *wrap.Particles, backendHandle *backend.Handle) Idle { 
  return Idle{ 
    stateId      : stateId,
    particles    : particles,
    backendHandle: backendHandle } 
}

func (idle Idle) Destroy() {}
func (idle Idle) PreEvent()    State { return idle }

func (idle Idle) OnTick()      State {
  var err error

  err = idle.particles.ReloadParticles(idle.backendHandle)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to reload particles: %s\n", err) 
  }

  err = object.AdvanceTick(idle.backendHandle, idle.stateId)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed advance tick after reloading particles: %s\n", err) 
  }

  return idle 
}

func (idle Idle) TickDelayed() bool  { 
  var (
    err       error 
    ticksLeft uint32 
  )

  ticksLeft, err = object.QueryTicksLeft(idle.backendHandle, idle.stateId)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to query ticks left: %s\n", err) 
  }

  return ticksLeft > 0 
}

func (idle Idle) OnTickDelay() State { 
  var err error 

  err = object.AdvanceTick(idle.backendHandle, idle.stateId)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed advance tick when delayed: %s\n", err) 
  }
  
  return idle 
}

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