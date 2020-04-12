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
}

func MakeIdle(backendHandle *backend.Handle) Idle {
  return Idle{ backendHandle: backendHandle }
}

func (idle Idle) OnTick() State {
  return idle
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
          var drawDot State 

          drawDot, err = MakeDrawDot(idle, idle.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})
          if err != nil {
            fmt.Fprintf(os.Stderr, "Could not draw dot at (%d,%d): %s\n", event.X, event.Y, err)
          } else {
            state = drawDot 
          }
        case sdl.BUTTON_RIGHT : 
          var drawLine State 

          drawLine, err = MakeDrawLine(idle, idle.backendHandle, backend.Position{X: event.X, Y: event.Y}, backend.Color{R: uint8(event.X%256), G: uint8((event.Y+70)%256), B: uint8((event.X+140)%256), A: 255})

          if err != nil {
            fmt.Fprintf(os.Stderr, "Could not start line at (%d,%d): %s\n", event.X, event.Y, err)
          } else {
            state = drawLine 
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

