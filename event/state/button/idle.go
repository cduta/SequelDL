package button

import (
	"os"
	"fmt"

  "../../../backend"
  "../../../sdlex"
  . "../../state"

  "github.com/veandco/go-sdl2/sdl"
)

type Idle struct {
	buttonEntityId  int64
  handle         *backend.Handle
}

func MakeIdle(buttonEntityId int64, handle *backend.Handle) (Idle, error) {
	var err error

	err = handle.ChangeState("button-idle", buttonEntityId)

  return Idle{ 
  	buttonEntityId : buttonEntityId, 
    handle         : handle }, err
}

func (idle Idle) Destroy() {}

func (idle Idle) OnTick() State { return idle }

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

func (idle Idle) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var (
    err       error
    state     State = idle
    collision bool 
  )

  switch event.State {
    case sdlex.BUTTON_PRESSED: 
      switch event.Button {
        case sdl.BUTTON_LEFT  : 
        	collision, err = idle.handle.HasEntityPixelCollision(idle.buttonEntityId, backend.Position{X: event.X, Y: event.Y})
        	if err != nil {
    				fmt.Fprintf(os.Stderr, "Failed to check if the button was hit: %s\n", err)
        		return idle  
        	}

        	if collision {
        		state, err = MakePressed(idle)
        		if err != nil {
    					fmt.Fprintf(os.Stderr, "Failed to create the button pressed state: %s\n", err)
    					state = idle
        		}
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

