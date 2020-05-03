package scene

import (
  "os"
  "fmt"

  "../../../backend"
  "../../../sdlex"
  . "../../state"

  "github.com/veandco/go-sdl2/sdl"
)

type Init struct {
  scene    *sdlex.Scene 
  images   *backend.Images
  renderer *sdl.Renderer
  idle      Idle 
}

func MakeInit(backendHandle *backend.Handle, scene *sdlex.Scene, renderer *sdl.Renderer) (Init, error) {
  var (
    err        error 
    images    *backend.Images
  )

  images, err = backendHandle.QueryImages(scene.Id) 

  return Init{ scene: scene, images: images, renderer: renderer, idle: MakeIdle(backendHandle, scene) }, err
}

func (init Init) Destroy() {
  if init.images != nil {
    init.images.Close()
  }
}

func (init Init) cancelInit() State {
  init.idle.Destroy()
  return MakeDone(init)
}

func (init Init) loadNextImage() (State, error) {
  var (
    err           error 
    backendImage *backend.Image
    sceneImage    sdlex.Image
    absolutePath  string
  )

  backendImage, err = init.images.Next() 
  if err != nil {
    return init.cancelInit(), err
  }

  if backendImage != nil {
    absolutePath, err = backend.ToAbsolutePath(backendImage.ImagePath)
    if err != nil {
      return init.cancelInit(), err
    }
    
    sceneImage, err = sdlex.MakeImage(init.renderer, backendImage.Id, backendImage.Name, absolutePath)
    if err != nil {
      return init.cancelInit(), err
    }

    init.scene.AddImage(&sceneImage)
    return init, err 
  } else {
    init.scene.IsReady(true)
    return Transition(init, init.idle), err 
  }
}

func (init Init) PreEvent() State { return init }

func (init Init) OnTick() State {
  var (
    err   error 
    state State
  )

  state, err = init.loadNextImage()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not load all images needed for the scene (%s): %s\n", init.scene.Name, err)
  }

  return state
}

func (init Init) OnQuit(event *sdl.QuitEvent) State { return MakeQuit(init) }

func (init Init) OnKeyboardEvent(event *sdl.KeyboardEvent) State {
  var state State = init

  switch event.State {
    case sdlex.BUTTON_PRESSED:  
      if event.Keysym.Mod & sdl.KMOD_CTRL > 0 {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(init)
        }
      } else {
        switch event.Keysym.Sym {
          case sdl.K_ESCAPE:
            state = MakeQuit(init)
        }
      }
  } 

  return state
}

func (init Init) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State { return init }
func (init Init) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State { return init }
func (init Init) PostEvent() State { return init }


