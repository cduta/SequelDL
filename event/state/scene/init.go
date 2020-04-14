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
}

func MakeInit(backendHandle *backend.Handle, scene *sdlex.Scene, renderer *sdl.Renderer) (Init, error) {
  var (
    err        error 
    images    *backend.Images
  )

  images, err = backendHandle.QueryImages(scene.Name) 

  return Init{ scene: scene, images: images, renderer: renderer }, err
}

func (init Init) Destroy() {
  if init.images != nil {
    init.images.Close()
  }
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
    return MakeEnd(init), err
  }

  if backendImage != nil {
    absolutePath, err = backend.ToAbsolutePath(backendImage.ImagePath)
    if err != nil {
      return MakeEnd(init), err
    }
    
    sceneImage, err = sdlex.MakeImage(init.renderer, backendImage.Id, backendImage.Name, absolutePath)
    if err != nil {
      return MakeEnd(init), err
    }

    init.scene.AddImage(&sceneImage)
    return init, err 
  } else {
    return MakeEnd(init), err 
  }
}

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

func (init Init) OnQuit(event *sdl.QuitEvent) State {
  return MakeQuit(init)
}

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

func (init Init) OnMouseMotionEvent(event *sdl.MouseMotionEvent) State {
  return init
}

func (init Init) OnMouseButtonEvent(event *sdl.MouseButtonEvent) State {
  return init
}

