package backend

func (handle *Handle) HasEntityPixelCollision(entityId int64, position Position) (bool, error) {
  var (
    err        error 
    row       *Row
    collision  bool = false 
  )

  row, err = handle.queryRow(`
SELECT TOTAL(? BETWEEN MAX(sc.x + en.x + hi.relative_x - sc.scroll_x            , sc.x) 
                   AND MIN(sc.x + en.x + hi.relative_x - sc.scroll_x + hi.width , sc.x + sc.width ) AND 
             ? BETWEEN MAX(sc.y + en.y + hi.relative_y - sc.scroll_y            , sc.y) 
                   AND MIN(sc.y + en.y + hi.relative_y - sc.scroll_y + hi.height, sc.y + sc.height)) > 0
FROM   scenes          AS sc,
       entities_scenes AS es,
       entities        AS en,
       hitboxes        AS hi  
WHERE  sc.id = es.scene_id 
AND    en.id = es.entity_id
AND    en.id = ?
AND    en.id = hi.entity_id
  `, position.X, position.Y, entityId)
  if err != nil {
    return collision, err
  }

  if row != nil {
    err = row.Scan(&collision)
  } 

  return collision, err 
}

/*
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
ORDER BY sp.entity_level, sp.sprite_level;
*/