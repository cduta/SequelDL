package object

import (
  "../../../backend"
)

func HasEntityPixelCollision(handle *backend.Handle, entityId int64, position backend.Position) (bool, error) {
  var (
    err        error 
    row       *backend.Row
    collision  bool = false 
  )

  row, err = handle.QueryRow(`
SELECT TOTAL(en.id = ?) AND TRUE
FROM (
  SELECT en.id
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
  ORDER BY en.level DESC, hi.level DESC, en.id DESC
  LIMIT 1
) AS en;
`, entityId, position.X, position.Y)
  if err != nil {
    return collision, err
  }

  if row != nil {
    err = row.Scan(&collision)
  } 

  return collision, err 
}
