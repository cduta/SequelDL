package state

import (
	"../../backend"
	"../../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type DrawLine struct {
	previousState  State
  backendHandle *backend.Handle
  sdlWrap       *sdlex.Wrap
  from           backend.Position
}

func MakeDrawLine(previousState State, backendHandle *backend.Handle, sdlWrap *sdlex.Wrap, from backend.Position) DrawLine {
	return DrawLine{
		previousState: previousState, 
		backendHandle: backendHandle, 
		sdlWrap      : sdlWrap, 
		from         : from}
}

func (drawLine DrawLine) OnQuit(event *sdl.QuitEvent) State {
	drawLine.sdlWrap.StopRunning()
	return MakeStop()
}

func (drawLine DrawLine) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
	var state State = drawLine

  switch event.Keysym.Sym {
    case sdl.K_ESCAPE:
      drawLine.sdlWrap.StopRunning()
      state = MakeStop()
  }

  return state
}

func (drawLine DrawLine) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
	switch event.State {
	  case sdl.BUTTON_LEFT :
    	// Preview line
	  case sdl.BUTTON_RIGHT:
	}

	return drawLine
}

func (drawLine DrawLine) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
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
		    	// Finish drawing line
		    default               :
		  }
  }

	return drawLine 
}