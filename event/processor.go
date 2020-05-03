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
  preProcess    bool 
  processEvents bool 
  processTicks  bool 
  postProcess   bool
  active        bool
}

func NewProcess(state state.State) Process {
  return Process{
    state        : state, 
    preProcess   : true,
    processEvents: true,  
    processTicks : true,
    postProcess  : true,
    active       : true}
}

func (process *Process) Destroy() {
  if process.state != nil {
    process.state.Destroy()
  }
}

func (process *Process) doProcessEvents(processEvents bool) {
  process.processEvents = processEvents
}

func (process *Process) doProcessTicks(processTicks bool) {
  process.processTicks = processTicks
}

func (process *Process) doPreProcess(preProcess bool) {
  process.preProcess = preProcess
}

func (process *Process) doPostProcess(postProcess bool) {
  process.postProcess = postProcess
}

func (process *Process) isActive(active bool) {
  process.active = active
}

type Processor struct {
  processes []Process
  sdlWrap    *sdlex.Wrap
  lastTick    time.Time
}

func (processor *Processor) Processes() []Process {
  return processor.processes
}

func NewProcessor(sdlWrap *sdlex.Wrap) *Processor {
  return &Processor{ processes: []Process{}, sdlWrap: sdlWrap, lastTick: time.Now() }
}

func (processor *Processor) AddProcess(process Process) {
  processor.processes = append(processor.processes, process)
}

func (processor *Processor) PreProcess() {
  var (
    i       int 
    process Process
  )
  
  for i, process = range processor.processes {
    if process.active && process.preProcess {
        processor.processes[i].state = process.state.PreEvent()      
    }
  }
}

func (processor *Processor) PostProcess() {
  var (
    i       int 
    process Process
  )

  for i, process = range processor.processes {
    if process.active && process.postProcess {
        processor.processes[i].state = process.state.PostEvent()      
    }
  }
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
      if process.active && process.processTicks {
        processor.processes[i].state = process.state.OnTick()
      }
    }
  }
}

func (processor *Processor) processEvents() {
  var (
    polledEvent sdl.Event
    i,j         int
    process     Process
  )

  for polledEvent = sdl.PollEvent(); polledEvent != nil; polledEvent = sdl.PollEvent() {
    for i, process = range processor.processes {

      if process.active && process.processEvents {
        switch event := polledEvent.(type) {
          case *sdl.QuitEvent:        processor.processes[i].state = process.state.OnQuit(event)
          case *sdl.KeyboardEvent:    processor.processes[i].state = process.state.OnKeyboardEvent(event)
          case *sdl.MouseMotionEvent: processor.processes[i].state = process.state.OnMouseMotionEvent(event)
          case *sdl.MouseButtonEvent: processor.processes[i].state = process.state.OnMouseButtonEvent(event)
        }

        switch processor.processes[i].state.(type) {
          case state.Quit: 
            processor.sdlWrap.StopRunning() 
            for j, _ = range processor.processes {
              processor.processes[j].isActive(false)
            }
          case state.Done:
            processor.processes[i].isActive(false)
        } 
      }
    }
  }
}

func (processor *Processor) removeInactiveProcesses() {
  var i int
  
  for i = 0; i < len(processor.processes); {
    if !processor.processes[i].active {
      processor.processes[i].Destroy()
      processor.processes[i] = processor.processes[ len(processor.processes)-1]
      processor.processes    = processor.processes[:len(processor.processes)-1]
    } else {
      i++
    }
  }  
}

func (processor *Processor) ProcessStates() {
  processor.PreProcess()
  processor.processTicks()
  processor.processEvents()
  processor.PostProcess()
  processor.removeInactiveProcesses()
}
