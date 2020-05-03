package backend

func (handle *Handle) HasEntityPixelCollision(entityId int64, position Position) (bool, error) {
  var (
    err        error 
    row       *Row
    collision  bool = false 
  )

  row, err = handle.queryRow(`
SELECT TOTAL(en.id = ?) AND TRUE
FROM   scenes          AS sc,
       entities_scenes AS es,
       entities        AS en,
       hitboxes        AS hi  
WHERE  sc.id = es.scene_id 
AND    en.id = es.entity_id
AND    en.id = hi.entity_id
AND    ? BETWEEN MAX( sc.x + en.x + hi.relative_x - sc.scroll_x            , sc.x             ) 
             AND MIN( sc.x + en.x + hi.relative_x - sc.scroll_x + hi.width , sc.x + sc.width  ) 
AND    ? BETWEEN MAX( sc.y + en.y + hi.relative_y - sc.scroll_y            , sc.y             ) 
             AND MIN( sc.y + en.y + hi.relative_y - sc.scroll_y + hi.height, sc.y + sc.height )
ORDER BY en.level DESC, hi.level DESC
LIMIT 1;
`, entityId, position.X, position.Y)
  if err != nil {
    return collision, err
  }

  if row != nil {
    err = row.Scan(&collision)
  } 

  return collision, err 
}
