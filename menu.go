package main

import (
  "os"
  "fmt"
  "./sdlex"
  "./backend"
  "./event"
  "./event/state"

  "github.com/veandco/go-sdl2/sdl"
)

func run() {
  var (
    err               error
    backendHandle     *backend.Handle
    wrap              *sdlex.Wrap
    processor         *event.Processor

    wrapArgs           sdlex.WrapArgs = sdlex.WrapArgs{ 
      DEFAULT_WINDOW_TITLE : "Menu Test", 
      DEFAULT_WINDOW_WIDTH : 1024, 
      DEFAULT_WINDOW_HEIGHT: 786,
      DEFAULT_FONT         : "DejaVuSansMono.ttf",
      DEFAULT_FONT_SIZE    : 30,
      DEFAULT_FPS          : 60,
     	DEFAULT_SHOW_FPS     : false}
  )

  backendHandle, err = backend.NewHandle()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize backend: %s\n", err)
    return 
  }
  defer backendHandle.Close()
  wrapArgs.Handle = backendHandle

  sdl.Init(sdl.INIT_EVERYTHING)
  defer sdl.Quit()

  wrap, err = sdlex.NewWrap(wrapArgs)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize SDL: %s\n", err)
    return
  }
  defer wrap.Quit()

  processor = event.NewProcessor(state.MakeIdle(backendHandle, wrap))

  for wrap.IsRunning() {
    wrap.PrepareFrame()
    processor.ProcessEvents()
    wrap.RenderFrame()
    wrap.ShowFrame()
  }

  sdl.Quit()
}

func main() {
  sdl.Main(run)
}