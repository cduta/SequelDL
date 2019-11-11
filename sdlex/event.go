package sdlex

import (

  "github.com/veandco/go-sdl2/sdl"
)

type EventHandlers struct {
  OnQuit             func(event *sdl.QuitEvent)
  OnKeyboardEvent    func(event *sdl.KeyboardEvent)
  OnMouseMotionEvent func(event *sdl.MouseMotionEvent)
  OnMouseButtonEvent func(event *sdl.MouseButtonEvent)
}

func (eventHandlers EventHandlers) ProcessEvents() {
  var polledEvent sdl.Event
  for polledEvent = sdl.PollEvent(); polledEvent != nil; polledEvent = sdl.PollEvent() {

    switch event := polledEvent.(type) {
      case *sdl.QuitEvent:
        eventHandlers.OnQuit(event)
      case *sdl.KeyboardEvent:
        eventHandlers.OnKeyboardEvent(event)
      case *sdl.MouseMotionEvent:
      	eventHandlers.OnMouseMotionEvent(event)
      case *sdl.MouseButtonEvent:
        eventHandlers.OnMouseButtonEvent(event)
    }
  }
}
