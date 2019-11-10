package main

import (
  "os"
  "fmt"
  "./sdlex"
  "./backend"

  "github.com/veandco/go-sdl2/sdl"
)

func run() int {
  var (
    err               error
    backendHandle     *backend.Handle
    sdlwrapArgs       sdlex.SDLWrapArgs = sdlex.SDLWrapArgs{ 
      DEFAULT_WINDOW_TITLE : "Menu Test", 
      DEFAULT_WINDOW_WIDTH : 1024, 
      DEFAULT_WINDOW_HEIGHT: 786,
      DEFAULT_FONT         : "DejaVuSansMono.ttf",
      DEFAULT_FONT_SIZE    : 30}
    sdlWrap           *sdlex.SDLWrap
    eventHandlers     sdlex.EventHandlers = sdlex.EventHandlers{
      OnQuit            : func(event *sdl.QuitEvent) { sdlWrap.StopRunning() },
      OnKeyboardEvent   : func(event *sdl.KeyboardEvent) {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            sdlWrap.StopRunning()
        }
      },
      OnMouseButtonEvent: func(event *sdl.MouseButtonEvent) {
        var err error
        /*
        var mouseButton, mouseState string

        switch event.Button {
          case sdl.BUTTON_LEFT : mouseButton = "Left" 
          case sdl.BUTTON_RIGHT: mouseButton = "Right"
          default              : mouseButton = "Unknown" 
        }

        switch event.State {
          case sdlex.BUTTON_PRESSED : mouseState = "pressed"
          case sdlex.BUTTON_RELEASED: mouseState = "released"
          default                   : mouseState = "unknown"
        }
        */

        if event.Button == sdl.BUTTON_LEFT && event.State == sdlex.BUTTON_PRESSED {
          //sdlWrap.RenderTestShape(event.X, event.Y)
          err = backendHandle.AddObject("line", event.X, event.Y)
          if err != nil {
            fmt.Fprintf(os.Stderr, "Could not place object at (%d,%d): %s\n", event.X, event.Y, err)
          }

        } 
      }}
  )

  backendHandle, err = backend.NewHandle()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize backend: %s\n", err)
    return 1
  }
  defer backendHandle.Close()
  sdlwrapArgs.Handle = backendHandle

  sdl.Init(sdl.INIT_EVERYTHING)
  defer sdl.Quit()

  sdlWrap, err = sdlex.NewSDLWrap(sdlwrapArgs)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize SDL: %s\n", err)
    return 1
  }
  defer sdlWrap.Quit()

  err = backendHandle.AddObject("line",100,100)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not create object: %s\n", err)
  }

  for sdlWrap.IsRunning() {
    sdlWrap.PrepareFrame()
    eventHandlers.ProcessEvents()
    sdlWrap.RenderFrame()
    sdlWrap.ShowFrame()
  }

  return 0
}

func main() {
  os.Exit(run())
}