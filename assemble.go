package assemble

import (
  "os"
  "fmt"
  "runtime/pprof"
  "flag"
  "./sdlex"
  "./backend"
  "./event"
)

type MakeProcessor func(backendHandle *backend.Handle, sdlWrap *sdlex.Wrap) (*event.Processor, error)

type settings struct {
  DEFAULT_SAVE_FILE_PATH string 
}

func parseArgs() settings {
  var (
    saveFilePathArg *string = flag.String("l", "", "load game")
    cpuprofileArg   *string = flag.String("c", "", "save cpu profile to file")

    saveFilePath string 
  )
  
  flag.Parse()
  
  if *saveFilePathArg != "" {
    saveFilePath = *saveFilePathArg
  }

  if *cpuprofileArg != "" {
    f, err := os.Create(*cpuprofileArg)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Failed to inizialize profiler: %s\n", err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }

  return settings{ DEFAULT_SAVE_FILE_PATH: saveFilePath }
}

func initialize(save backend.Save, load backend.Load, render sdlex.Rendering) (*sdlex.Wrap, *backend.Handle, error) {
  var (
    err            error
    settings       settings = parseArgs()
    backendHandle *backend.Handle
    sdlWrap       *sdlex.Wrap
  )

  backendHandle, err = backend.MakeHandle(save, load, settings.DEFAULT_SAVE_FILE_PATH)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize backend: %s\n", err)
    return nil, nil, err
  }

  sdlWrap, err = sdlex.MakeWrap(backendHandle, render)
  if err != nil {
    backendHandle.Close()
    fmt.Fprintf(os.Stderr, "Failed to inizialize SDL wrapper: %s\n", err)
    return nil, nil, err
  }

  return sdlWrap, backendHandle, err
}

func mainLoop(sdlWrap *sdlex.Wrap, eventProcessor *event.Processor) {
  for sdlWrap.IsRunning() {
    eventProcessor.ProcessStates()
    sdlWrap.PrepareFrame()
    sdlWrap.RenderFrame()
    sdlWrap.ShowFrame()
  }

  return
}

func Run(save backend.Save, load backend.Load, makeProcessor MakeProcessor, render sdlex.Rendering) {
  var (
    err            error
    backendHandle *backend.Handle
    sdlWrap       *sdlex.Wrap
    processor     *event.Processor
  )

  sdlWrap, backendHandle, err = initialize(save, load, render)
  if err != nil {
    return
  }
  defer backendHandle.Close()
  defer sdlWrap.Quit()

  processor, err = makeProcessor(backendHandle, sdlWrap)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize event processor: %s\n", err)
    return
  }

  mainLoop(sdlWrap, processor)
}

