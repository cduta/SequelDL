package state

import (
  "fmt"
  "os"
  "database/sql"

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

	_, err = addDot(backendHandle, position, color)

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

func (drawDot DrawDot) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  var err error

  if sdlex.IsMouseMotionState(event.State, sdl.BUTTON_LEFT) {
      _, err = addDot(drawDot.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
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

func addDot(handle *backend.Handle, pos backend.Position, color backend.Color) (int64, error) {
  var (
    err          error
    result       sql.Result
    lastInsertId int64
  ) 

  result, err = handle.Exec(`
BEGIN IMMEDIATE;
INSERT OR ROLLBACK INTO objects DEFAULT VALUES;
`)

  if err != nil {
    return lastInsertId, err 
  }

  lastInsertId, err = result.LastInsertId()

  if err != nil {
    return lastInsertId, err
  }

  result, err = handle.Exec(`
INSERT OR ROLLBACK INTO dots(object_id, x, y) VALUES (?, ?, ?);
INSERT OR ROLLBACK INTO colors(object_id, r, g, b, a) VALUES (?, ?, ?, ?, ?);
COMMIT;
`, lastInsertId, pos.X, pos.Y, lastInsertId, color.R, color.G, color.B, color.A)
  
  return lastInsertId, err
}
          