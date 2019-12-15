package state

import (
  "fmt"
  "os"

  "../../backend"
  "../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type DrawLine struct {
  previousState  State
  backendHandle *backend.Handle
  objectId       int64
}

func MakeDrawLine(previousState State, backendHandle *backend.Handle, here backend.Position, color backend.Color) (DrawLine, error) {
  var (
    err          error
    lastInsertId int64
  ) 

  lastInsertId, err = backend.InsertLine(backendHandle, here, here, color)

  if err != nil {
    return DrawLine{}, err
  }

  return DrawLine{
    previousState: previousState, 
    backendHandle: backendHandle, 
    objectId     : lastInsertId}, 
    err
}

func (drawLine DrawLine) OnQuit(event *sdl.QuitEvent) State {
  return MakeQuit()
}

func (drawLine DrawLine) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var state State = drawLine

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      switch event.Keysym.Sym {
        case sdl.K_ESCAPE:
          state = MakeQuit()
      }
  }

  return state
}

func (drawLine DrawLine) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  var ( 
    err   error 
    state State = drawLine
  )
 
  if sdlex.IsMouseMotionState(event.State, sdl.BUTTON_RIGHT) {
    err = backend.UpdateLineThere(drawLine.backendHandle, drawLine.objectId, backend.Position{X: event.X, Y: event.Y})
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not update line coordinates: %s\n", err)
    }
  }

  return state
}

func (drawLine DrawLine) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var state State = drawLine

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
        case sdl.BUTTON_RIGHT : 
          state = drawLine.previousState 
        default               :
      }
  }

  return state 
}