package idle

import (
  "../../../../sdlex"
  . "../../../../event/state"

  "github.com/veandco/go-sdl2/sdl"
)

type Idle struct {}

func MakeIdle() Idle { return Idle{} }

func (idle Idle) Destroy() {}
func (idle Idle) PreEvent()    State { return idle }
func (idle Idle) OnTick()      State { return idle }
func (idle Idle) TickDelayed() bool  { return false }
func (idle Idle) OnTickDelay() State { return idle }
func (idle Idle) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(idle) }
func (idle Idle) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var (
    state State = idle
  )

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      if event.Keysym.Mod & sdl.KMOD_CTRL > 0 {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(idle)
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