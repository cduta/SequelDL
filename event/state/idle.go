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
  wrap          *sdlex.Wrap
}

func MakeIdle(backendHandle *backend.Handle, wrap *sdlex.Wrap) Idle {
	return Idle{
		backendHandle: backendHandle,
		wrap         : wrap}
}

func (idle Idle) OnQuit(event *sdl.QuitEvent) State {
	idle.wrap.StopRunning()

	return idle
}

func (idle Idle) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  switch event.Keysym.Sym {
    case sdl.K_ESCAPE:
      idle.wrap.StopRunning()
  }

  return idle
}

func (idle Idle) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
	var err error 

	if event.State == sdl.BUTTON_LEFT {
    err = idle.backendHandle.AddDot(backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not place object at (%d,%d): %s\n", event.X, event.Y, err)
    }      		
	}

	return idle
}

func (idle Idle) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  var err error
  /*
  var mouseButton, mouseState string

  switch event.Button {
    case sdl.BUTTON_LEFT : mouseButton = "Left" 
    case sdl.BUTTON_RIGHT: mouseButton = "Right"
    default              : mouseButton = "Unknown" 
  }

  switch event.State {
    case sdlex.BUTTON_PRESSED : mouseState = "pressed"
    case sdlex.BUTTON_RELEASED: mouseState = "released"
    default                   : mouseState = "unknown"
  }
  */

  if event.Button == sdl.BUTTON_LEFT && event.State == sdlex.BUTTON_PRESSED {
    err = idle.backendHandle.AddDot(backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not place object at (%d,%d): %s\n", event.X, event.Y, err)
    }

  } 

	return idle 
}