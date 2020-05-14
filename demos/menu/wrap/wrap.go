package wrap

import (
  "../../../backend"
	"../../../sdlex"
)

type MenuWrap struct {
  scene *Scene
}

func MakeMenuWrap() *MenuWrap {
	return &MenuWrap{ scene: nil }
}

func (menuWrap *MenuWrap) Destroy() {}

func (menuWrap *MenuWrap) Initialize(sclWrap *sdlex.SdlWrap, handle *backend.Handle) error {
  var (
    err        error
    menuScene *Scene
  )

  menuScene, err = MakeScene("menu", handle)
  if err != nil {
  	return err 
  }

  menuWrap.scene = menuScene

  return err 
}

func (menuWrap *MenuWrap) IsReady() bool {
  return menuWrap.scene.IsReady()
}

func (menuWrap *MenuWrap) Render(sdlWrap *sdlex.SdlWrap, handle *backend.Handle) error {
  var err error

  err = menuWrap.RenderDots(sdlWrap, handle)
  err = menuWrap.RenderLines(sdlWrap, handle)
  err = menuWrap.RenderScene(sdlWrap, handle)

  return err
}

func (menuWrap *MenuWrap) Scene() *Scene {
	return menuWrap.scene
} 