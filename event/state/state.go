package state

import (
  "github.com/veandco/go-sdl2/sdl"
)

type State interface {
	Destroy()
	PreEvent() State
	OnTick() State
	TickDelayed() bool
	OnTickDelay() State
  OnQuit(event *sdl.QuitEvent) State
  OnKeyboardEvent(event *sdl.KeyboardEvent) State
  OnMouseMotionEvent(event *sdl.MouseMotionEvent) State
  OnMouseButtonEvent(event *sdl.MouseButtonEvent) State
  PostEvent() State
}

func Transition(old State, new State) State {
	old.Destroy()
	return new
}