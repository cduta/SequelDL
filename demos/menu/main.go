package main

import (
  "os"
  "fmt"
	"../.."
  "../../sdlex"
  "../../backend"
  "../../event"
  "./state/draw"
  "./state/scene"
  "./state/button"
)

func MakeMenuProcessor(backendHandle *backend.Handle, sdlWrap *sdlex.Wrap) (*event.Processor, error) {
  var (
    err             error
    eventProcessor *event.Processor
    init            scene.Init
    initialScene   *sdlex.Scene 
    idle            button.Idle
  )

  eventProcessor = event.NewProcessor(sdlWrap)

  eventProcessor.AddProcess(event.NewProcess(draw.MakeIdle(backendHandle)))

  initialScene, err = sdlex.MakeScene("menu", backendHandle)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to make initial scene: %s\n", err)
    return eventProcessor, err
  }

  init, err = scene.MakeInit(sdlWrap.Handle(), initialScene, sdlWrap.Renderer())
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to make initializing scene state: %s\n", err)
    return eventProcessor, err
  }
  sdlWrap.SetScene(initialScene)
  eventProcessor.AddProcess(event.NewProcess(init))

  idle, err = button.MakeIdle(1, backendHandle)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to make initial idle state: %s\n", err) 
  }
  eventProcessor.AddProcess(event.NewProcess(idle))

  idle, err = button.MakeIdle(2, backendHandle)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to make initial idle state: %s\n", err) 
  }
  eventProcessor.AddProcess(event.NewProcess(idle))

  return eventProcessor, err
}

func main() {
  assemble.Run(MakeMenuProcessor)
}