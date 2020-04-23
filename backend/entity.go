package backend

func (handle *Handle) HasEntityPixelCollision(entityId int64, position Position) (bool, error) {
  var (
    err        error 
    row       *Row
    collision  bool = false 
  )

  row, err = handle.queryRow(`
SELECT TOTAL(? BETWEEN en.x + hi.relative_x AND en.x + hi.relative_x + hi.width AND 
             ? BETWEEN en.y + hi.relative_y AND en.y + hi.relative_y + hi.height   ) > 0
FROM   entities AS en,
       hitboxes AS hi  
WHERE  en.id = ?
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
