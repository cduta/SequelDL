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

  rows, err = handle.QueryRows(`
SELECT sp.id, 
       sp.name,
       sp.sprite_x      + MAX(sp.scene_x - sp.sprite_x, 0)                                          AS x,
       sp.sprite_y      + MAX(sp.scene_y - sp.sprite_y, 0)                                          AS y,
       sp.sprite_width  - MAX(sp.scene_x - sp.sprite_x, 0) 
                        - MAX((sp.sprite_x + sp.sprite_width)  - (sp.scene_x + sp.scene_width) , 0) AS width, 
       sp.sprite_height - MAX(sp.scene_y - sp.sprite_y, 0) 
                        - MAX((sp.sprite_y + sp.sprite_height) - (sp.scene_y + sp.scene_height), 0) AS height,
                          MAX(sp.scene_x - sp.sprite_x, 0)                                          AS clip_x,
                          MAX(sp.scene_y - sp.sprite_y, 0)                                          AS clip_y
FROM (
  SELECT sp.id, 
         sp.name, 
         sc.x + en.x + sp.relative_x - sc.scroll_x AS sprite_x,      -- ⎫  
         sc.y + en.y + sp.relative_y - sc.scroll_y AS sprite_y,      -- ⎬ Absolute sprite position and size on screen 
         sp.width                                  AS sprite_width,  -- ⎪  
         sp.height                                 AS sprite_height, -- ⎭
         sp.level                                  AS sprite_level,
         en.level                                  AS entity_level, 
         sc.x                                      AS scene_x,       -- ⎫  
         sc.y                                      AS scene_y,       -- ⎬ Absolute scene position and size on screen
         sc.width                                  AS scene_width,   -- ⎪  
         sc.height                                 AS scene_height   -- ⎭
  FROM   sprites         AS sp,
         states_sprites  AS ss,
         states          AS st,
         entities        AS en,
         entities_scenes AS es,
         scenes          AS sc
  WHERE  en.visible
  AND    sp.id       = ss.sprite_id
  AND    ss.state_id = st.id 
  AND    st.id       = en.state_id
  AND    en.id       = es.entity_id
  AND    es.scene_id = sc.id 
  AND    sc.id = ?
) AS sp
WHERE sp.sprite_x                    < sp.scene_x + sp.scene_width
AND   sp.sprite_y                    < sp.scene_y + sp.scene_height
AND   sp.sprite_x + sp.sprite_width  > sp.scene_x 
AND   sp.sprite_y + sp.sprite_height > sp.scene_y
ORDER BY sp.entity_level, sp.sprite_level, sp.id;
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
    err            error
    spriteId       int64 
    name           string 
    x, y, w, h     int64
    clip_x, clip_y int64
  )

  if !sprites.Rows.Next() {
    return nil, err
  }

  err = sprites.Rows.Scan(&spriteId, &name, &x, &y, &w, &h, &clip_x, &clip_y)  

  return &Sprite{ 
    Id        : spriteId, 
    Name      : name, 
    SrcLayout : &sdl.Rect{X: int32(clip_x) , Y: int32(clip_y) , W: int32(w), H: int32(h)}, 
    DestLayout: &sdl.Rect{X: int32(x)      , Y: int32(y)      , W: int32(w), H: int32(h)} }, err
}
