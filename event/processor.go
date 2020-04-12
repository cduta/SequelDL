package event

import (
  "./state"
  "../sdlex"

  "github.com/veandco/go-sdl2/sdl"

  "time"
  "fmt"
  "os"
)

const TICK_INTERVAL_IN_MILLISECONDS int64 = 20

type Processor struct {
  states   []state.State
  sdlWrap  *sdlex.Wrap
  lastTick  time.Time
}

func NewProcessor(initialState state.State, sdlWrap *sdlex.Wrap) *Processor {
  return &Processor{ states: []state.State{initialState}, sdlWrap: sdlWrap, lastTick: time.Now() }
}

func (processor *Processor) processTicks() {
  var ( 
    err             error 
    now             time.Time     = time.Now()
    msSinceLastTick int64         = now.Sub(processor.lastTick).Milliseconds()
    ticks           int64         = msSinceLastTick/TICK_INTERVAL_IN_MILLISECONDS
    restString      string        = fmt.Sprintf("%dms", msSinceLastTick%TICK_INTERVAL_IN_MILLISECONDS)
    rest            time.Duration
    t               int64
    i               int 
    state           state.State
  )

  rest, err = time.ParseDuration(restString)
  if err != nil {    
    fmt.Fprintf(os.Stderr, "Rest duration string for tick was malformed: %s\n", restString)
  }

  processor.lastTick = now.Add(-rest)

  for t = 0; t < ticks; t++ {
    for i, state = range processor.states {
      processor.states[i] = state.OnTick()
    }
  }
}

func (processor *Processor) processEvents() {
  var (
    polledEvent sdl.Event
    i           int
    s           state.State
  )

  for polledEvent = sdl.PollEvent(); polledEvent != nil; polledEvent = sdl.PollEvent() {
    for i, s = range processor.states {

      switch event := polledEvent.(type) {
        case *sdl.QuitEvent:        processor.states[i] = s.OnQuit(event)
        case *sdl.KeyboardEvent:    processor.states[i] = s.OnKeyboardEvent(event)
        case *sdl.MouseMotionEvent: processor.states[i] = s.OnMouseMotionEvent(event)
        case *sdl.MouseButtonEvent: processor.states[i] = s.OnMouseButtonEvent(event)
      }

      switch processor.states[i].(type) {
        case state.Quit: processor.sdlWrap.StopRunning() 
      } 
    }
  }
}

func (processor *Processor) ProcessStates() {
  processor.processTicks()
  processor.processEvents()
}
