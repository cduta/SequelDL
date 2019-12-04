package state

import (
  "fmt"
  "os"
  "database/sql"

  "../../backend"
  "../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type DrawLine struct {
  previousState  State
  backendHandle *backend.Handle
  objectId       int64
}

func MakeDrawLine(previousState State, backendHandle *backend.Handle, from backend.Position, color backend.Color) (DrawLine, error) {
  var (
    err          error
    result       sql.Result
    lastInsertId int64
  ) 

  result, err = backendHandle.Exec(`
BEGIN IMMEDIATE;
INSERT OR ROLLBACK INTO objects DEFAULT VALUES;
`)

  if err != nil {
    return DrawLine{}, err 
  }

  lastInsertId, err = result.LastInsertId()

  if err != nil {
    return DrawLine{}, err
  }

  _, err = backendHandle.Exec(`
INSERT OR ROLLBACK INTO lines(object_id, here_x, here_y, there_x, there_y) VALUES (?, ?, ?, ?, ?);
INSERT OR ROLLBACK INTO colors(object_id, r, g, b, a) VALUES (?, ?, ?, ?, ?);
COMMIT;
`, lastInsertId, from.X, from.Y, from.X, from.Y, lastInsertId, color.R, color.G, color.B, color.A)

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
    _, err = drawLine.backendHandle.Exec(`
  UPDATE lines   
  SET    there_x = ?, there_y = ? 
  WHERE  object_id = ?
    `, event.X, event.Y, drawLine.objectId)
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