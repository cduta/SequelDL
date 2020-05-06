package button

import (
  "fmt"
  "os"

  "../../../../backend"
  "../../../../sdlex"
  . "../../../../event/state"

  "github.com/veandco/go-sdl2/sdl"
)

type Pressed struct {
	idle            Idle 
  buttonEntityId  int64
  handle         *backend.Handle
}

func MakePressed(idle Idle) (Pressed, error) {
	var err error

	err = idle.handle.ChangeState("button-pressed", idle.buttonEntityId)

  return Pressed{ 
  	idle           : idle, 
    buttonEntityId : idle.buttonEntityId,
    handle         : idle.handle }, err
}

func (pressed Pressed) Destroy() {}
func (pressed Pressed) PreEvent()    State { return pressed }
func (pressed Pressed) OnTick()      State { return pressed }
func (pressed Pressed) TickDelayed() bool  { return false }
func (pressed Pressed) OnTickDelay() State { return pressed }

func (pressed Pressed) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(pressed) }

func (pressed Pressed) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var (
    state State = pressed
  )

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      if event.Keysym.Mod & sdl.KMOD_CTRL > 0 {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(pressed)
        }
      } else {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(pressed)
        }
      }
  } 

  return state
}

func (pressed Pressed) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State { return pressed }

func (pressed Pressed) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var (
    err   error
    state State = pressed
  )

  switch event.State {
    case sdlex.BUTTON_PRESSED: 
      switch event.Button {
        case sdl.BUTTON_LEFT  : 
        case sdl.BUTTON_RIGHT : 
        default               : 
      }
    case sdlex.BUTTON_RELEASED: 
      switch event.Button {
        case sdl.BUTTON_LEFT  :  
            err = pressed.handle.ChangeState("button-idle", pressed.idle.buttonEntityId)
            if err != nil {
              fmt.Fprintf(os.Stderr, "Failed to release button: %s\n", err)
              return pressed
            }
            state = pressed.idle
        case sdl.BUTTON_RIGHT : 
        default               :
      }
  }

  return state 
}

func (pressed Pressed) PostEvent() State { return pressed }
