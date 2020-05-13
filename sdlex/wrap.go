package sdlex

import (
  "os"
  "fmt"

  "../backend"

  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/img"
  "github.com/veandco/go-sdl2/ttf"
  "github.com/veandco/go-sdl2/gfx"
)

type sdlWrapArgs struct {
  DEFAULT_WINDOW_TITLE   string
  DEFAULT_WINDOW_WIDTH   int32
  DEFAULT_WINDOW_HEIGHT  int32 
  DEFAULT_FONT           string
  DEFAULT_FONT_SIZE      int
  DEFAULT_FPS            uint32 
  DEFAULT_SHOW_FPS       bool
  Handle                *backend.Handle
}

type Wrap interface {
  Destroy()
  IsReady() bool
  Render(sdlWrap *SdlWrap) error
}

type SdlWrap struct {
  running     bool
  showFPS     bool
  window     *sdl.Window
  renderer   *sdl.Renderer
  font       *ttf.Font
  fpsManager *gfx.FPSmanager
  handle     *backend.Handle
  Scene      *Scene 
}

func MakeSdlWrap(backendHandle *backend.Handle) (*SdlWrap, error) {
  var (
    err         error
    options     backend.Options
    args        sdlWrapArgs
    window     *sdl.Window
    renderer   *sdl.Renderer
    font       *ttf.Font
    fpsManager *gfx.FPSmanager = new(gfx.FPSmanager)
  )

  args.Handle = backendHandle

  options, err = backendHandle.QueryOptions()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to query options: %s\n", err)
    return nil, err
  }

  args = sdlWrapArgs{ 
    DEFAULT_WINDOW_TITLE :        options.WindowTitle,       
    DEFAULT_WINDOW_WIDTH :  int32(options.WindowWidth),     
    DEFAULT_WINDOW_HEIGHT:  int32(options.WindowHeight),    
    DEFAULT_FONT         :        options.DefaultFont,      
    DEFAULT_FONT_SIZE    :    int(options.DefaultFontSize), 
    DEFAULT_FPS          : uint32(options.FPS),             
    DEFAULT_SHOW_FPS     :        options.ShowFPS,          
    Handle               :        backendHandle}

  sdl.Init(sdl.INIT_EVERYTHING)

  if args.Handle == nil {
    return nil, fmt.Errorf("Backend handle not defined")
  }

  window, err = sdl.CreateWindow(args.DEFAULT_WINDOW_TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
    args.DEFAULT_WINDOW_WIDTH, args.DEFAULT_WINDOW_HEIGHT, sdl.WINDOW_SHOWN)
  if err != nil {
    return nil, err
  }

  renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
  if err != nil {
    return nil, err
  }

  err = img.Init(img.INIT_PNG)
  if err != nil {
    return nil, err
  }

  err = ttf.Init();
  if err != nil {
    return nil, err
  }

  font, err = ttf.OpenFont(args.DEFAULT_FONT, args.DEFAULT_FONT_SIZE);
  if err != nil {
    return nil, err
  }

  gfx.InitFramerate(fpsManager)
  gfx.SetFramerate(fpsManager, args.DEFAULT_FPS)

  return &SdlWrap{
    running   : true,
    showFPS   : args.DEFAULT_SHOW_FPS,
    window    : window,      
    renderer  : renderer,
    font      : font,
    fpsManager: fpsManager,
    handle    : args.Handle}, err
}

func (sdlWrap SdlWrap) Quit() {
  sdlWrap.StopRunning()
  sdlWrap.window.Destroy()
  sdlWrap.renderer.Destroy()
  sdlWrap.font.Close()
  img.Quit()  
  sdl.Quit()
}

func (sdlWrap SdlWrap) Renderer() *sdl.Renderer {
  return sdlWrap.renderer
}

func (sdlWrap SdlWrap) Handle() *backend.Handle {
  return sdlWrap.handle
}

func (sdlWrap SdlWrap) IsRunning() bool {
  return sdlWrap.running
}

func (sdlWrap *SdlWrap) StopRunning() {
  sdlWrap.running = false
}

func (sdlWrap *SdlWrap) SetScene(scene *Scene) {
  sdlWrap.Scene = scene
}

func (sdlWrap SdlWrap) PrepareFrame() {
  gfx.FramerateDelay(sdlWrap.fpsManager)
  sdlWrap.renderer.SetDrawColor(0, 0, 0, 255)
  sdlWrap.renderer.Clear()
}

func (sdlWrap *SdlWrap) RenderGenerics() {
  var err error

  if sdlWrap.showFPS {
    err = sdlWrap.RenderFramerate(0,0)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Failed to render FPS: %s\n", err)
    }
  }
}

func (sdlWrap *SdlWrap) RenderWrap(wrap Wrap) {
  var err error 

  if wrap != nil && wrap.IsReady() {
    err = wrap.Render(sdlWrap)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Failed to render objects: %s\n", err)
    }
  }
}

func (sdlWrap SdlWrap) ShowFrame() {
  sdlWrap.renderer.Present()
}

func (sdlWrap SdlWrap) newTextTexture(str string) (*sdl.Texture, int32, int32, error) {
  var (
    err      error
    surface *sdl.Surface
    texture *sdl.Texture
  )

  surface, err = sdlWrap.font.RenderUTF8Blended(str, sdl.Color{255, 0, 0, 255})
  if err != nil {
    return texture, 0, 0, err
  }
  defer surface.Free()

  texture, err = sdlWrap.renderer.CreateTextureFromSurface(surface)
  return texture, surface.W, surface.H, err
}

func (sdlWrap SdlWrap) renderText(text string, x, y int32) (error) {
  var (
    err      error
    w, h     int32
    texture *sdl.Texture
  )

  if texture, w, h, err = sdlWrap.newTextTexture(text); err != nil {
    return err
  }
  defer texture.Destroy()

  return sdlWrap.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
}

func (sdlWrap SdlWrap) RenderFramerate(x, y int32) (error) {
  var (
    success bool
    framerate int
    framerateString string
  )

  framerate, success = gfx.GetFramerate(sdlWrap.fpsManager)

  if success {
    framerateString = fmt.Sprintf("%d FPS ", framerate)
  } else {
    framerateString = "N/A FPS"
  } 

  return sdlWrap.renderText(framerateString, x, y)
}

func PrintRendererInfos() error {
  var (
    err      error 
    drivers  int 
    info    *sdl.RendererInfo = new(sdl.RendererInfo)
  )

  drivers, err = sdl.GetNumRenderDrivers()
  if err != nil {
    return err
  }

  for i := 0; i < drivers; i++ {
    _, err = sdl.GetRenderDriverInfo(i, info)
    if err != nil {
      return err
    }
    fmt.Printf("Render Driver Index:%d\n%+v\n", i, info)
  }

  return err
}