package draw

import (
  "fmt"
  "os"

  "../../../backend"
  "../../../sdlex"
  . "../../state"

  "github.com/veandco/go-sdl2/sdl"
)

type Dot struct {
  previousState  State
  backendHandle *backend.Handle
}

func MakeDot(previousState State, backendHandle *backend.Handle, position backend.Position, color backend.Color) (Dot, error) {
  var err error

	_, err = backend.InsertDot(backendHandle, position, color)

  if err != nil {
  	return Dot{}, err
  }       

  return Dot{
    previousState: previousState, 
    backendHandle: backendHandle}, 
    err
}

func (dot Dot) Destroy() {}

func (dot Dot) PreEvent() State { return dot }
func (dot Dot) OnTick()   State { return dot }
func (dot Dot) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(dot) }

func (dot Dot) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var state State = dot

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      switch event.Keysym.Sym {
        case sdl.K_ESCAPE:
          state = MakeQuit(dot)
      }
  }

  return state
}

func (dot Dot) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  var err error

  if sdlex.IsMouseMotionState(event.State, sdl.BUTTON_LEFT) {
      _, err = backend.InsertDot(dot.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
      if err != nil {
        fmt.Fprintf(os.Stderr, "Could not draw dot at (%d,%d): %s\n", event.X, event.Y, err)
      }       
  }

  return dot
}

func (dot Dot) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var state State = dot

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
          state = dot.previousState 
        case sdl.BUTTON_RIGHT : 
        default               :
      }
  }

  return state 
}
          
func (dot Dot) PostEvent() State { return dot }