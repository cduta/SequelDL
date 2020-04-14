package state

import (
  "github.com/veandco/go-sdl2/sdl"
)

type End struct {
	state State
}

func MakeEnd(state State) End {
  return End{state: state}
}

func (end End) Destroy() {
	end.state.Destroy()
}

func (end End) OnTick() State {
	return end
}

func (end End) OnQuit(event *sdl.QuitEvent) State {
  return end
}

func (end End) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  return end
}

func (end End) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  return end
}

func (end End) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  return end 
}