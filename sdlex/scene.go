package sdlex

type Scene struct {
	Name			 string
  Images     map[int64]*Image
  Ready      bool
}

func MakeScene(name string) Scene {
	return Scene{ Name: name, Images: make(map[int64]*Image), Ready: false }
}

func (scene *Scene) Destroy() {
	var image *Image

	for _, image = range scene.Images {
		image.Destroy()
	}
	scene.Images = make(map[int64]*Image)
	scene.Ready = false
}

func (scene *Scene) AddImage(image *Image) {
	scene.Images[image.Id] = image
}

func (scene *Scene) IsSceneReady(ready bool) {
	scene.Ready = ready
}