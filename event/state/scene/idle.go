package scene

import (
  "os"
  "fmt"

  "../../../backend"
  "../../../sdlex"
  . "../../state"

  "github.com/veandco/go-sdl2/sdl"
)

const (
  SCROLL_FIRST_DELAY  = 5
  SCROLL_REPEAT_DELAY = 2
)

type Idle struct {
  backendHandle    *backend.Handle
  scene            *sdlex.Scene 
  scrollX           int32
  scrollY           int32
  firstScroll       bool
  delayFirstRepeat  int
  delayRepeats      int  
}

func MakeIdle(backendHandle *backend.Handle, scene *sdlex.Scene) Idle {
  var idle Idle = Idle{ 
    backendHandle       : backendHandle, 
    scene               : scene}

  idle.reset()

  return idle
}

func (idle *Idle) reset() {
  idle.scrollX          = 0
  idle.scrollY          = 0
  idle.firstScroll      = true 
  idle.delayFirstRepeat = SCROLL_FIRST_DELAY
  idle.delayRepeats     = 0
}

func (idle Idle) scrolling() bool {
  return idle.scrollX != 0 || idle.scrollY != 0
}

func (idle Idle) Destroy() {}

func (idle Idle) PreEvent() State { return idle }

func (idle Idle) OnTick()   State { 
  var err error 

  if idle.scrolling() {
    err = idle.backendHandle.ScrollScene(idle.scene.Id, idle.scrollX, idle.scrollY)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Could not scroll scene (%s) in (%d, %d) direction: %s\n", idle.scene.Name, idle.scrollX, idle.scrollY, err)
    }

    idle.firstScroll = false

    if idle.delayFirstRepeat == 0 {
      idle.delayRepeats = SCROLL_REPEAT_DELAY
    }
  }

  return idle  
}

func (idle Idle) TickDelayed() bool  { 
  return !idle.firstScroll && (idle.delayFirstRepeat > 0 || idle.delayRepeats > 0) 
}

func (idle Idle) OnTickDelay() State { 
  if idle.delayFirstRepeat > 0 {
    idle.delayFirstRepeat--
  }

  if idle.delayRepeats > 0 {
    idle.delayRepeats--
  }

  return idle 
}

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
          case (event.Keysym.Sym == sdl.K_UP    || 
                event.Keysym.Sym == sdl.K_LEFT  || 
                event.Keysym.Sym == sdl.K_RIGHT || 
                event.Keysym.Sym == sdl.K_DOWN) &&
               !sdlex.IsRepeatButtonPress(event) :
            if  event.Keysym.Sym == sdl.K_UP {
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
    case sdlex.BUTTON_RELEASED:  
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
               event.Keysym.Sym == sdl.K_DOWN :
            if event.Keysym.Sym == sdl.K_UP {
              idle.scrollY +=  1
            }
            if event.Keysym.Sym == sdl.K_LEFT {
              idle.scrollX +=  1
            }
            if event.Keysym.Sym == sdl.K_RIGHT {
              idle.scrollX += -1
            }
            if event.Keysym.Sym == sdl.K_DOWN {
              idle.scrollY += -1
            }

            if (!idle.scrolling()) {
              idle.reset()
            }
            state = idle
        }
      }
  } 

  return state
}

func (idle Idle) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State { return idle }
func (idle Idle) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State { return idle }
func (idle Idle) PostEvent() State { return idle }
