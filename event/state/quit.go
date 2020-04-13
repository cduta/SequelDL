package state

import (
  "github.com/veandco/go-sdl2/sdl"
)

type Quit struct {
	state State
}

func MakeQuit(state State) Quit {
  return Quit{state: state}
}

func (quit Quit) Destroy() {
	quit.state.Destroy()
}

func (quit Quit) OnTick() State {
	return quit
}

func (quit Quit) OnQuit(event *sdl.QuitEvent) State {
  return quit
}

func (quit Quit) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  return quit
}

func (quit Quit) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  return quit
}

func (quit Quit) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  return quit 
}