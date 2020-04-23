package backend

import (
  "github.com/veandco/go-sdl2/sdl"
)

type Sprites struct {
  *Rows
}

type Sprite struct {
  Id          int64
  Name        string
  SrcLayout  *sdl.Rect
  DestLayout *sdl.Rect
}

func (handle *Handle) QuerySprites(sceneId int64) (*Sprites, error) {
  var (
    err   error 
    rows *Rows
  )

  rows, err = handle.queryRows(`
SELECT sp.id, 
       sp.name, 
       en.x + sp.relative_x, 
       en.y + sp.relative_y, 
       sp.width, 
       sp.height
FROM   sprites         AS sp,
       states_sprites  AS ss,
       states          AS st,
       entities        AS en,
       entities_scenes AS es,
       scenes          AS sc
WHERE  sp.id       = ss.sprite_id
AND    ss.state_id = st.id 
AND    st.id       = en.state_id
AND    en.id       = es.entity_id
AND    es.scene_id = ?
`, sceneId)
  if rows == nil || err != nil {
    return nil, err 
  }

  return &Sprites{ Rows: rows }, err
}

func (sprites Sprites) Close() {
  sprites.Rows.Close()
}

func (sprites Sprites) Next() (*Sprite, error) {
  var (
    err        error
    spriteId   int64 
    name       string 
    x, y, w, h int64
  )

  if !sprites.Rows.next() {
    return nil, err
  }

  err = sprites.Rows.Scan(&spriteId, &name, &x, &y, &w, &h)  

  return &Sprite{ 
    Id        : spriteId, 
    Name      : name, 
    SrcLayout : &sdl.Rect{X: 0, Y: 0, W: int32(w), H: int32(h)}, 
    DestLayout: &sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)} }, err
}
