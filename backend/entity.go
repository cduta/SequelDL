package backend

import (
  "database/sql"
)

func (handle *Handle) HasEntityPixelCollision(entityId int64, position Position) (bool, error) {
  var (
    err        error 
    row       *sql.Row
    collision  bool 
  )

  row, err = handle.queryRow(`
SELECT TOTAL(? BETWEEN en.x + hi.relative_x AND en.x + hi.relative_x + hi.width AND 
             ? BETWEEN en.y + hi.relative_y AND en.y + hi.relative_y + hi.height   ) > 0
FROM   entities AS en,
       hitboxes AS hi  
WHERE  en.id = ?
AND    en.id = hi.entity_id
  `, position.X, position.Y, entityId)
  if err != nil && err != sql.ErrNoRows {
    return collision, err
  }

  if err != sql.ErrNoRows {
    err = row.Scan(&collision)
  } else {
    collision = false
  }

  return collision, err 
}
