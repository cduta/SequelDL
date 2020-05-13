package backend

type Images struct {
  *Rows
}

type Image struct {
  Id        int64
  Name      string
  ImagePath string
}

func (handle *Handle) QueryImages(sceneId int64) (*Images, error) {
  var (
    err   error 
    rows *Rows
  )

  rows, err = handle.QueryRows(`
SELECT DISTINCT im.id, im.name, im.image_path 
FROM (
    SELECT DISTINCT ss.animation_group
    FROM   states_sprites  AS ss,
           states          AS st,
           entities        AS en,
           entities_scenes AS es,
           scenes          AS sc
    WHERE  ss.state_id = st.id 
    AND    st.id       = en.state_id
    AND    en.id       = es.entity_id
    AND    es.scene_id = ?
  )              AS _,
  states_sprites AS ss,
  sprites        AS sp,
  images         AS im 
WHERE _.animation_group = ss.animation_group
AND   ss.sprite_id      = sp.id
AND   sp.image_id       = im.id
`, sceneId)
  if err != nil {
    return nil, err 
  }

  return &Images{ Rows: rows }, err
}

func (images Images) Close() {
  images.Rows.Close()
}

func (images Images) Next() (*Image, error) {
  var (
    err       error
    imageId   int64
    name      string 
    imagePath string 
  )

  if !images.Rows.Next() {
    return nil, err
  }

  err = images.Rows.Scan(&imageId, &name, &imagePath)  

  return &Image{ Id: imageId, Name: name, ImagePath: imagePath }, err
}
