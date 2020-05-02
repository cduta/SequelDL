package state

import (
  "github.com/veandco/go-sdl2/sdl"
)

type Done struct {
	state State
}

func MakeDone(state State) Done {
  return Done{state: state}
}

func (done Done) Destroy() {
	done.state.Destroy()
}

func (done Done) OnTick() State {
	return done
}

func (done Done) OnQuit(event *sdl.QuitEvent) State {
  return done
}

func (done Done) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  return done
}

func (done Done) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  return done
}

func (done Done) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  return done 
}