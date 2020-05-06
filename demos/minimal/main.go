package main

import (
	"../.."
  "../../sdlex"
  "../../backend"
  "../../event"
  "./state/idle"
)

func MakeMinimalProcessor(_ *backend.Handle, sdlWrap *sdlex.Wrap) (*event.Processor, error) {
  var (
    err             error
    eventProcessor *event.Processor
  )
  
  eventProcessor = event.NewProcessor(sdlWrap)
  eventProcessor.AddProcess(event.NewProcess(idle.MakeIdle()))

  return eventProcessor, err
}

func main() {
  assemble.Run(MakeMinimalProcessor)
}