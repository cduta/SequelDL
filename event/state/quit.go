package state

import (
  "github.com/veandco/go-sdl2/sdl"
)

type Quit struct {}

func MakeQuit() Quit {
	return Quit{}
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