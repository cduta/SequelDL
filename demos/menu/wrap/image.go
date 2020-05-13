package wrap

import (
  "github.com/veandco/go-sdl2/sdl"
  "github.com/veandco/go-sdl2/img"
)

type Image struct {
	Id       int64
  Name     string
  Path     string
  texture *sdl.Texture
}

func MakeImage(renderer *sdl.Renderer, id int64, name string, absolutePath string) (Image, error) {
	var (
		err      error 
		texture *sdl.Texture
	)

	texture, err = img.LoadTexture(renderer, absolutePath)
  return Image{ Id: id, Name: name, Path: absolutePath, texture: texture }, err
}

func (image *Image) Destroy() {
	if image.texture != nil {
		image.texture.Destroy()
		image.texture = nil
	}
}