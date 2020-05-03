package draw

import (
  "fmt"
  "os"

  "../../../backend"
  "../../../sdlex"
  . "../../state"

  "github.com/veandco/go-sdl2/sdl"
)

type Line struct {
  previousState  State
  backendHandle *backend.Handle
  objectId       int64
}

func MakeLine(previousState State, backendHandle *backend.Handle, here backend.Position, color backend.Color) (Line, error) {
  var (
    err          error
    lastInsertId int64
  ) 

  lastInsertId, err = backend.InsertLine(backendHandle, here, here, color)

  if err != nil {
    return Line{}, err
  }

  return Line{
    previousState: previousState, 
    backendHandle: backendHandle, 
    objectId     : lastInsertId}, 
    err
}

func (line Line) Destroy() {}
func (line Line) PreEvent()    State { return line }
func (line Line) OnTick()      State { return line }
func (line Line) TickDelayed() bool  { return false }
func (line Line) OnTickDelay() State { return line }
func (line Line) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(line) }

func (line Line) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var state State = line

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      switch event.Keysym.Sym {
        case sdl.K_ESCAPE:
          state = MakeQuit(line)
      }
  }

  return state
}

func (line Line) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  var ( 
    err   error 
    state State = line
  )
 
  if sdlex.IsMouseMotionState(event.State, sdl.BUTTON_RIGHT) {
    err = backend.UpdateLineThere(line.backendHandle, line.objectId, backend.Position{X: event.X, Y: event.Y})
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not update line coordinates: %s\n", err)
    }
  }

  return state
}

func (line Line) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var state State = line

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
          state = line.previousState 
        default               :
      }
  }

  return state 
}

func (line Line) PostEvent() State { return line }