package sdlex

import (
  "os"
  "fmt"
  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/ttf"
  "github.com/veandco/go-sdl2/gfx"
  "../backend"
)

type SDLWrapArgs struct {
  DEFAULT_WINDOW_TITLE   string
  DEFAULT_WINDOW_WIDTH   int32
  DEFAULT_WINDOW_HEIGHT  int32 
  DEFAULT_FONT           string
  DEFAULT_FONT_SIZE      int
  DEFAULT_FPS            uint32 
  DEFAULT_SHOW_FPS       bool
  Handle                *backend.Handle
}

type SDLWrap struct {
  running     bool
  showFPS     bool
  window     *sdl.Window
  renderer   *sdl.Renderer
  font       *ttf.Font
  fpsManager *gfx.FPSmanager
  handle     *backend.Handle
}

func NewSDLWrap(args SDLWrapArgs) (*SDLWrap, error) {
  var (
    err         error
    window     *sdl.Window
    renderer   *sdl.Renderer
    font       *ttf.Font
    fpsManager *gfx.FPSmanager = new(gfx.FPSmanager)
  )

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

  return &SDLWrap{
    running   : true,
    showFPS   : args.DEFAULT_SHOW_FPS,
    window    : window,      
    renderer  : renderer,
    font      : font,
    fpsManager: fpsManager,
    handle    : args.Handle}, err
}

func (sdlWrap SDLWrap) Quit() {
  sdlWrap.StopRunning()
  sdlWrap.window.Destroy()
  sdlWrap.renderer.Destroy()
  sdlWrap.font.Close()
}

func (sdlWrap SDLWrap) IsRunning() bool {
  return sdlWrap.running
}

func (sdlWrap *SDLWrap) StopRunning() {
  sdlWrap.running = false
}

func (sdlWrap SDLWrap) PrepareFrame() {
  sdlWrap.renderer.SetDrawColor(0, 0, 0, 255)
  sdlWrap.renderer.Clear()
}

func (sdlWrap SDLWrap) RenderFrame() {
  var err error

  gfx.FramerateDelay(sdlWrap.fpsManager)

  if sdlWrap.showFPS {
    err = sdlWrap.RenderFramerate(0,0)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Failed to render FPS: %s\n", err)
    }
  }

  err = sdlWrap.renderObjects()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to render lines: %s\n", err)
  }
}


func (sdlWrap SDLWrap) ShowFrame() {
  sdlWrap.renderer.Present()
}

func (sdlWrap SDLWrap) newTextTexture(str string) (*sdl.Texture, int32, int32, error) {
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

func (sdlWrap SDLWrap) renderText(text string, x, y int32) (error) {
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

func (sdlWrap SDLWrap) RenderFramerate(x, y int32) (error) {
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

func (sdlWrap SDLWrap) RenderDot(x, y int32) {
  sdlWrap.renderer.SetDrawColor(143, 143, 143, 255)
  sdlWrap.renderer.DrawPoint(x,y)  
}

func (sdlWrap SDLWrap) renderDots() error {
  var (
    err   error 
    dots *backend.Dots
    dot  *backend.Dot
  )

  dots, err = sdlWrap.handle.QueryDots()
  if err != nil {
    return err
  }
  defer dots.Close()

  for dot, err = dots.Next(); err == nil && dot != nil; dot, err = dots.Next() {
    sdlWrap.RenderDot(dot.X,dot.Y)
  }

  return err
}

func (sdlWrap SDLWrap) renderObjects() error {
  var err error

  err = sdlWrap.renderDots()

  return err
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