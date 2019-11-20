package state

import (
  "os"
  "fmt"

	"../../backend"
	"../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type Idle struct {
  backendHandle *backend.Handle
  sdlWrap       *sdlex.Wrap
}

func MakeIdle(backendHandle *backend.Handle, sdlWrap *sdlex.Wrap) Idle {
	return Idle{
		backendHandle: backendHandle,
		sdlWrap      : sdlWrap}
}

func (idle Idle) OnQuit(event *sdl.QuitEvent) State {
	idle.sdlWrap.StopRunning()
	return MakeStop()
}

func (idle Idle) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
	var state State = idle

  switch event.Keysym.Sym {
    case sdl.K_ESCAPE:
      idle.sdlWrap.StopRunning()
      state = MakeStop()
  }

  return state
}

func (idle Idle) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
	var err error 

	switch event.State {
	  case sdl.BUTTON_LEFT :
    	err = idle.backendHandle.AddDot(backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
	    if err != nil {
	      fmt.Fprintf(os.Stderr, "Could not place object at (%d,%d): %s\n", event.X, event.Y, err)
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
    	    err = idle.backendHandle.AddDot(backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
			    if err != nil {
			      fmt.Fprintf(os.Stderr, "Could not place object at (%d,%d): %s\n", event.X, event.Y, err)
			    }
		    case sdl.BUTTON_RIGHT : state = MakeDrawLine(idle, idle.backendHandle, idle.sdlWrap, backend.Position{X: event.X, Y: event.Y})
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