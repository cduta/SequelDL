package scene

import (
  "os"
  "fmt"

  "../../../backend"
  "../../../sdlex"
  . "../../state"

  "github.com/veandco/go-sdl2/sdl"
)

type Idle struct {
  backendHandle *backend.Handle
  scene         *sdlex.Scene 
  scrollX        int32
  scrollY        int32
}

func MakeIdle(backendHandle *backend.Handle, scene *sdlex.Scene) Idle {
  return Idle{ backendHandle: backendHandle, scene: scene, scrollX: 0, scrollY: 0 }
}

func (idle Idle) Destroy() {}
func (idle Idle) PreEvent() State { 
  idle.scrollX = 0
  idle.scrollY = 0

  return idle 
}
func (idle Idle) OnTick()   State { return idle }
func (idle Idle) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(idle) }

func (idle Idle) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var state State = idle

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      if event.Keysym.Mod & sdl.KMOD_CTRL > 0 {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(idle)
        }
      } else {
        switch {
          case event.Keysym.Sym == sdl.K_ESCAPE:
            state = MakeQuit(idle)
          case event.Keysym.Sym == sdl.K_UP    || 
               event.Keysym.Sym == sdl.K_LEFT  || 
               event.Keysym.Sym == sdl.K_RIGHT || 
               event.Keysym.Sym == sdl.K_DOWN  :
            if event.Keysym.Sym == sdl.K_UP {
              idle.scrollY += -1
            }
            if event.Keysym.Sym == sdl.K_LEFT {
              idle.scrollX += -1
            }
            if event.Keysym.Sym == sdl.K_RIGHT {
              idle.scrollX +=  1
            }
            if event.Keysym.Sym == sdl.K_DOWN {
              idle.scrollY +=  1
            }
            state = idle
        }
      }
  } 

  return state
}

func (idle Idle) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State { return idle }
func (idle Idle) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State { return idle }
func (idle Idle) PostEvent() State { 
  var err error 

  if idle.scrollX != 0 || idle.scrollY != 0 {
    err = idle.backendHandle.ScrollScene(idle.scene.Id, idle.scrollX, idle.scrollY)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not scroll scene (%s) in (%d, %d) direction: %s\n", idle.scene.Name, idle.scrollX, idle.scrollY, err)
    }
  }

  return idle 
}
