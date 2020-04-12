package state

import (
  "fmt"
  "os"

  "../../backend"
  "../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type DrawDot struct {
  previousState  State
  backendHandle *backend.Handle
}

func MakeDrawDot(previousState State, backendHandle *backend.Handle, position backend.Position, color backend.Color) (DrawDot, error) {
  var err error

	_, err = backend.InsertDot(backendHandle, position, color)

  if err != nil {
  	return DrawDot{}, err
  }       

  return DrawDot{
    previousState: previousState, 
    backendHandle: backendHandle}, 
    err
}

func (drawDot DrawDot) OnQuit(event *sdl.QuitEvent) State {
  return MakeQuit()
}

func (drawDot DrawDot) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var state State = drawDot

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      switch event.Keysym.Sym {
        case sdl.K_ESCAPE:
          state = MakeQuit()
      }
  }

  return state
}

func (drawDot DrawDot) OnTick() State {
  return drawDot
}

func (drawDot DrawDot) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  var err error

  if sdlex.IsMouseMotionState(event.State, sdl.BUTTON_LEFT) {
      _, err = backend.InsertDot(drawDot.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
      if err != nil {
        fmt.Fprintf(os.Stderr, "Could not draw dot at (%d,%d): %s\n", event.X, event.Y, err)
      }       
  }

  return drawDot
}

func (drawDot DrawDot) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var state State = drawDot

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
          state = drawDot.previousState 
        case sdl.BUTTON_RIGHT : 
        default               :
      }
  }

  return state 
}
          