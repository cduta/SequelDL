package event

import (
	"./state"
	"../sdlex"

  "github.com/veandco/go-sdl2/sdl"
)

type Processor struct {
	state    state.State       
	sdlWrap *sdlex.Wrap
}

func NewProcessor(initialState state.State, sdlWrap *sdlex.Wrap) *Processor {
	return &Processor{ state: initialState, sdlWrap: sdlWrap }
}

func (processor *Processor) ProcessEvents() {
  var polledEvent sdl.Event

  for polledEvent = sdl.PollEvent(); polledEvent != nil; polledEvent = sdl.PollEvent() {

    switch event := polledEvent.(type) {
      case *sdl.QuitEvent:        processor.state = processor.state.OnQuit(event)
      case *sdl.KeyboardEvent:    processor.state = processor.state.OnKeyboardEvent(event)
      case *sdl.MouseMotionEvent: processor.state = processor.state.OnMouseMotionEvent(event)
      case *sdl.MouseButtonEvent: processor.state	= processor.state.OnMouseButtonEvent(event)
    }

    switch processor.state.(type) {
    	case state.Quit: processor.sdlWrap.StopRunning() 
    } 
  }
}
