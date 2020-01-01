package main

import (
  "os"
  "fmt"
  "runtime/pprof"
  "flag"
  "./sdlex"
  "./backend"
  "./event"
  "./event/state"

  "github.com/veandco/go-sdl2/sdl"
)

func run() {
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

  eventProcessor = event.NewProcessor(state.MakeIdle(backendHandle), sdlWrap)

  for sdlWrap.IsRunning() {
    sdlWrap.PrepareFrame()
    eventProcessor.ProcessEvents()
    sdlWrap.RenderFrame()
    sdlWrap.ShowFrame()
  }

  sdl.Quit()
}

func main() {
  var cpuprofile *string = flag.String("c", "", "write cpu profile to file")
  flag.Parse()
  if *cpuprofile != "" {
      f, err := os.Create(*cpuprofile)
      if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to inizialize profiler: %s\n", err)
      }
      pprof.StartCPUProfile(f)
      defer pprof.StopCPUProfile()
  }
  run()
}