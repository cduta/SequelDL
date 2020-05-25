package object

import (
  "../../../backend"
)

func QueryTicksLeft(handle *backend.Handle, stateId int64) (uint32, error) {
  var (
    err        error 
    row       *backend.Row
    ticksLeft  int64 
  )

  row, err = handle.QueryRow(`
SELECT s.ticks_left
FROM   states AS s
WHERE  s.ticks_left IS NOT NULL 
AND    s.id = ?;
`, stateId)
  if err != nil {
    return 0, err
  }

  if row != nil {
    err = row.Scan(&ticksLeft)
  } 

  return uint32(ticksLeft), err 
}

func AdvanceTick(handle *backend.Handle, stateId int64) error {
  var err error

  _, err = handle.Exec(`
BEGIN IMMEDIATE;
UPDATE OR ROLLBACK states
SET   ticks_left = 
  CASE 
    WHEN ticks_left > 0 THEN ticks_left - 1
    ELSE ticks 
  END
WHERE id = ?;
COMMIT;
`, stateId)

  return err
}