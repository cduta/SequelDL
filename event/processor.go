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

type Process struct {
  state         state.State
  processEvents bool 
  processTicks  bool 
  active        bool
}

func NewProcess(state state.State) Process {
  return Process{
    state        : state, 
    processEvents: true,  
    processTicks : true,
    active       : true}
}

func (process *Process) doProcessEvents(processEvents bool) {
  process.processEvents = processEvents
}

func (process *Process) doProcessTicks(processTicks bool) {
  process.processTicks = processTicks
}

type Processor struct {
  processes []Process
  sdlWrap    *sdlex.Wrap
  lastTick    time.Time
}

func NewProcessor(initialState state.State, sdlWrap *sdlex.Wrap) *Processor {
  return &Processor{ processes: []Process{NewProcess(initialState)}, sdlWrap: sdlWrap, lastTick: time.Now() }
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
    process         Process
  )

  rest, err = time.ParseDuration(restString)
  if err != nil {    
    fmt.Fprintf(os.Stderr, "Rest duration string for tick was malformed: %s\n", restString)
  }

  processor.lastTick = now.Add(-rest)

  for t = 0; t < ticks; t++ {
    for i, process = range processor.processes {
      processor.processes[i].state = process.state.OnTick()
    }
  }
}

func (processor *Processor) processEvents() {
  var (
    polledEvent sdl.Event
    i           int
    process     Process
  )

  for polledEvent = sdl.PollEvent(); polledEvent != nil; polledEvent = sdl.PollEvent() {
    for i, process = range processor.processes {

      switch event := polledEvent.(type) {
        case *sdl.QuitEvent:        processor.processes[i].state = process.state.OnQuit(event)
        case *sdl.KeyboardEvent:    processor.processes[i].state = process.state.OnKeyboardEvent(event)
        case *sdl.MouseMotionEvent: processor.processes[i].state = process.state.OnMouseMotionEvent(event)
        case *sdl.MouseButtonEvent: processor.processes[i].state = process.state.OnMouseButtonEvent(event)
      }

      switch processor.processes[i].state.(type) {
        case state.Quit: 
          processor.sdlWrap.StopRunning() 
          processor.processes[i].active = false
      } 
    }
  }
}



func (processor *Processor) removeDefunctProcesses() {
  var i int
  
  for i = 0; i < len(processor.processes); {
    if !processor.processes[i].active {
      processor.processes[i] = processor.processes[len(processor.processes)-1]
      processor.processes = processor.processes[:len(processor.processes)-1]
    } else {
      i++
    }
  }  
}

func (processor *Processor) ProcessStates() {
  processor.processTicks()
  processor.processEvents()
  processor.removeDefunctProcesses()
}
