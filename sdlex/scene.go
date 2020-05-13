package sdlex

import (
  "../backend"
)

type Scene struct {
	Id     int64
	Name   string
  Images map[int64]*Image
  Ready  bool
}

func MakeScene(sceneName string, backendHandle *backend.Handle) (*Scene, error) {
	var (
		err     error 
		sceneId int64 = 1
	)

	sceneId, err = backendHandle.QuerySceneId(sceneName)
	if err != nil {
		return nil, err
	}

	return &Scene{ Id: sceneId, Name: sceneName, Images: make(map[int64]*Image), Ready: false }, err
}

func (scene *Scene) Destroy() {
	var image *Image

	for _, image = range scene.Images {
		image.Destroy()
	}
	scene.Images = make(map[int64]*Image)
	scene.Ready = false
}

func (scene *Scene) Clear() {
	scene.Destroy()
}

func (scene *Scene) AddImage(image *Image) {
	scene.Images[image.Id] = image
}

func (scene *Scene) SetReady(ready bool) {
	scene.Ready = ready
}

func (scene *Scene) IsReady() bool {
	return scene != nil && scene.Ready
}

func (sdlWrap SdlWrap) RenderSprite(sprite *backend.Sprite) {
	sdlWrap.renderer.Copy(sdlWrap.Scene.Images[sprite.Id].texture, sprite.SrcLayout, sprite.DestLayout)
}

func (sdlWrap SdlWrap) RenderScene() error {
  var (
    err      error 
    sprites *backend.Sprites
    sprite  *backend.Sprite
  )

  if !sdlWrap.Scene.IsReady() {
  	return err
  }

  sprites, err = sdlWrap.handle.QuerySprites(sdlWrap.Scene.Id)
  if sprites == nil || err != nil {
    return err
  }
  defer sprites.Close()

  for sprite, err = sprites.Next(); err == nil && sprite != nil; sprite, err = sprites.Next() {
    sdlWrap.RenderSprite(sprite)
  }

  return err
}