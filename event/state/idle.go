package state

import (
  "os"
  "fmt"
  "database/sql"

	"../../backend"
	"../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type Idle struct {
  backendHandle *backend.Handle
}

func MakeIdle(backendHandle *backend.Handle) Idle {
	return Idle{ backendHandle: backendHandle }
}

func (idle Idle) OnQuit(event *sdl.QuitEvent) State {
	return MakeQuit()
}

func (idle Idle) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
	var (
		err   error 
		state State = idle
	)

  switch event.State {
  	case sdlex.BUTTON_PRESSED:  
			if event.Keysym.Mod & sdl.KMOD_CTRL > 0 {
		  	switch event.Keysym.Sym {
			    case sdl.K_ESCAPE:
			      state = MakeQuit()
			    case sdl.K_s:
			    	err = idle.backendHandle.Save("save.db")
			    	if err != nil {
	      			fmt.Fprintf(os.Stderr, "Could not save: %s\n", err)
	      		}
			  }
			} else {
		  	switch event.Keysym.Sym {
			    case sdl.K_ESCAPE:
			      state = MakeQuit()
			  }
			}
	} 

  return state
}

func (idle Idle) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
	var err error 

	switch event.State {
	  case sdl.BUTTON_LEFT :
    	_, err = addDot(idle.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
	    if err != nil {
	      fmt.Fprintf(os.Stderr, "Could not place dot at (%d,%d): %s\n", event.X, event.Y, err)
	    }      	
	  case sdl.BUTTON_RIGHT:
	}

	return idle
}

func (idle Idle) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var (
  	err   error
  	state State = idle
  )

  switch event.State {
    case sdlex.BUTTON_PRESSED: 
  	  switch event.Button {
		    case sdl.BUTTON_LEFT  : 
    	    _, err = addDot(idle.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
			    if err != nil {
			      fmt.Fprintf(os.Stderr, "Could not place dot at (%d,%d): %s\n", event.X, event.Y, err)
			    }
		    case sdl.BUTTON_RIGHT : 
		    	var drawLine State 
		    	drawLine, err = MakeDrawLine(idle, idle.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})

		    	if err == nil {
		    		state = drawLine 
		    	} else {
	      		fmt.Fprintf(os.Stderr, "Could not start line at (%d,%d): %s\n", event.X, event.Y, err)
		    	}
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