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

type MakeProcessor func(backendHandle *backend.Handle, sdlWrap *sdlex.SdlWrap, wrap sdlex.Wrap) (*event.Processor, error)

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
  }

  return settings{ DEFAULT_SAVE_FILE_PATH: saveFilePath }
}

func initialize(save backend.Save, load backend.Load) (*sdlex.SdlWrap, *backend.Handle, error) {
  var (
    err            error
    settings       settings = parseArgs()
    backendHandle *backend.Handle
    sdlWrap       *sdlex.SdlWrap
  )

  backendHandle, err = backend.MakeHandle(save, load, settings.DEFAULT_SAVE_FILE_PATH)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize backend: %s\n", err)
    return nil, nil, err
  }

  sdlWrap, err = sdlex.MakeSdlWrap(backendHandle)
  if err != nil {
    backendHandle.Close()
    fmt.Fprintf(os.Stderr, "Failed to inizialize SDL wrapper: %s\n", err)
    return nil, nil, err
  }

  return sdlWrap, backendHandle, err
}

func mainLoop(sdlWrap *sdlex.SdlWrap, eventProcessor *event.Processor, wrap sdlex.Wrap) {
  for sdlWrap.IsRunning() {
    eventProcessor.ProcessStates()
    sdlWrap.PrepareFrame()
    sdlWrap.RenderGenerics()
    sdlWrap.RenderWrap(wrap)
    sdlWrap.ShowFrame()
  }

  return
}

func Run(save backend.Save, load backend.Load, makeProcessor MakeProcessor, wrap sdlex.Wrap) {
  var (
    err            error
    backendHandle *backend.Handle
    sdlWrap       *sdlex.SdlWrap
    processor     *event.Processor
  )

  if wrap == nil {
    fmt.Fprintf(os.Stderr, "When calling run, you did not provide a wrap. Aborting.")
  }

  sdlWrap, backendHandle, err = initialize(save, load)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize SDL wrapper or backend handler: %s\n", err)
    return
  }
  defer backendHandle.Close()
  defer sdlWrap.Quit()

  err = wrap.Initialize(sdlWrap, backendHandle)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize wrapper: %s\n", err)
    return
  }
  defer wrap.Destroy()

  processor, err = makeProcessor(backendHandle, sdlWrap, wrap)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize event processor: %s\n", err)
    return
  }

  mainLoop(sdlWrap, processor, wrap)
  pprof.StopCPUProfile()
}

