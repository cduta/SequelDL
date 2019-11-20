package state

import (
  "github.com/veandco/go-sdl2/sdl"
)

type Stop struct {}

func MakeStop() Stop {
	return Stop{}
}

func (stop Stop) OnQuit(event *sdl.QuitEvent) State {
	return stop
}

func (stop Stop) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  return stop
}

func (stop Stop) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
	return stop
}

func (stop Stop) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
	return stop 
}