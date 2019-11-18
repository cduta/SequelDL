package state

import (
  "github.com/veandco/go-sdl2/sdl"
)

type State interface {
  OnQuit(event *sdl.QuitEvent) State
  OnKeyboardEvent(event *sdl.KeyboardEvent) State
  OnMouseMotionEvent(event *sdl.MouseMotionEvent) State
  OnMouseButtonEvent(event *sdl.MouseButtonEvent) State
}
