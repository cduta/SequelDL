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
  "./event/state/scene"
  "./event/state/button"
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

func InitializeEventProcessor(backendHandle *backend.Handle, sdlWrap *sdlex.Wrap) (*event.Processor, error) {
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

func run(settings settings) {
  var (
    err             error
    backendHandle  *backend.Handle
    sdlWrap        *sdlex.Wrap
    eventProcessor *event.Processor
    options         backend.Options
    sdlWrapArgs     sdlex.WrapArgs
  )

  backendHandle, err = backend.NewHandle(settings.DEFAULT_SAVE_FILE_PATH)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize backend: %s\n", err)
    return 
  }
  defer backendHandle.Close()

  sdlWrapArgs.Handle = backendHandle

  options, err = backendHandle.QueryOptions()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to query options: %s\n", err)
    return     
  }

  sdlWrapArgs = sdlex.WrapArgs{ 
    DEFAULT_WINDOW_TITLE :        options.WindowTitle,       
    DEFAULT_WINDOW_WIDTH :  int32(options.WindowWidth),     
    DEFAULT_WINDOW_HEIGHT:  int32(options.WindowHeight),    
    DEFAULT_FONT         :        options.DefaultFont,      
    DEFAULT_FONT_SIZE    :    int(options.DefaultFontSize), 
    DEFAULT_FPS          : uint32(options.FPS),             
    DEFAULT_SHOW_FPS     :        options.ShowFPS,          
    Handle               :        backendHandle}

  sdl.Init(sdl.INIT_EVERYTHING)
  defer sdl.Quit()

  sdlWrap, err = sdlex.NewWrap(sdlWrapArgs)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to inizialize SDL: %s\n", err)
    return
  }
  defer sdlWrap.Quit()

  eventProcessor, err = InitializeEventProcessor(backendHandle, sdlWrap)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to initialize event processor: %s\n", err)
    return
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