package backend

type Images struct {
  *Objects
}

type Image struct {
  Id        int64
  Name      string
  ImagePath string
}

func (handle Handle) QueryImages(sceneName string) (*Images, error) {
  var (
    err      error 
    objects *Objects
  )

  objects, err = handle.queryObjects(`
SELECT i.id, i.name, i.image_path
FROM   images        AS i, 
       images_scenes AS isc,
       scenes        AS s
WHERE  isc.image_id  = i.id
AND    s.id         = isc.scene_id
AND    s.name       = ?;
`, sceneName)
  if err != nil {
    return nil, err 
  }

  return &Images{ Objects: objects }, err
}

func (images Images) Close() {
  images.Objects.Close()
}

func (images Images) Next() (*Image, error) {
  var (
    err       error
    imageId   int64
    name      string 
    imagePath string 
  )

  if !images.Objects.next() {
    return nil, err
  }

  err = images.Objects.rows.Scan(&imageId, &name, &imagePath)  

  return &Image{ Id: imageId, Name: name, ImagePath: imagePath }, err
}
