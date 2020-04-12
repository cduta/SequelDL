package main

import (
  "os"
  "fmt"
  "runtime/pprof"
  "flag"
  "./sdlex"
  "./backend"
  "./event"
  "./event/state/draw"

  "github.com/veandco/go-sdl2/sdl"
)

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

func run(settings settings) {
  var (
    err             error
    backendHandle  *backend.Handle
    sdlWrap        *sdlex.Wrap
    eventProcessor *event.Processor

    sdlWrapArgs     sdlex.WrapArgs = sdlex.WrapArgs{ 
      DEFAULT_WINDOW_TITLE : "Menu Test", 
      DEFAULT_WINDOW_WIDTH : 1024, 
      DEFAULT_WINDOW_HEIGHT: 786,
      DEFAULT_FONT         : "font/DejaVuSansMono.ttf",
      DEFAULT_FONT_SIZE    : 30,
      DEFAULT_FPS          : 60,
      DEFAULT_SHOW_FPS     : true}
  )

  backendHandle, err = backend.NewHandle()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize backend: %s\n", err)
    return 
  }
  defer backendHandle.Close()

  sdlWrapArgs.Handle = backendHandle

  sdl.Init(sdl.INIT_EVERYTHING)
  defer sdl.Quit()

  sdlWrap, err = sdlex.NewWrap(sdlWrapArgs)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize SDL: %s\n", err)
    return
  }
  defer sdlWrap.Quit()

  eventProcessor = event.NewProcessor(draw.MakeIdle(backendHandle), sdlWrap)

  if settings.DEFAULT_SAVE_FILE_PATH != "" {
    backendHandle.Load(settings.DEFAULT_SAVE_FILE_PATH)
  }

  for sdlWrap.IsRunning() {
    eventProcessor.ProcessStates()
    sdlWrap.PrepareFrame()
    sdlWrap.RenderFrame()
    sdlWrap.ShowFrame()
  }

  sdl.Quit()
}

func main() {
  run(parseArgs())
}