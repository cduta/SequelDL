package object

import (  
  "../../../backend"

  "database/sql"
)

func IgnitePaper(backendHandle *backend.Handle, position backend.Position) error {
  var err error

  _, err = backendHandle.Exec(`
BEGIN IMMEDIATE;
UPDATE OR ROLLBACK entities_states 
SET   state_id = 2
WHERE state_id = 1
AND   EXISTS (SELECT 1 
              FROM   entities AS e 
              WHERE  e.id       = entity_id
              AND    (e.x, e.y) = (?,?));
COMMIT;
`, position.X, position.Y)

  return err 
}

func BurnPaper(backendHandle *backend.Handle) (int64, error) {
  var (
  	err          error
  	result       sql.Result
  	extinguished int64
  	ignited      int64
  )

  result, err = backendHandle.Exec(`
BEGIN IMMEDIATE;

UPDATE OR ROLLBACK entities_states 
SET   state_id  = 3
WHERE state_id  = 2
AND   abs(random() % 101) <= (SELECT io.value
                              FROM   integer_options AS io 
                              WHERE  io.id = 5);
`)
	if err != nil {
		return 0, err
	}

	extinguished, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

result, err = backendHandle.Exec(`
UPDATE OR ROLLBACK entities_states
SET   state_id  = 2
WHERE state_id  = 1
AND   EXISTS (SELECT 1
							FROM   entities         AS e1,
							       entities_states  AS es,
							       states_particles AS sp,
							       particles        AS p,
							       entities         AS e2
							WHERE  es.entity_id   = e1.id 
							AND    es.state_id    = 2
							AND    es.state_id    = sp.state_id 
							AND    sp.particle_id = p.id  
							AND    e2.id          = entities_states.entity_id
							AND    (e2.x, e2.y)   = (e1.x + p.relative_x, e1.y + p.relative_y))
AND   abs(random() % 101) <= (SELECT io.value
                              FROM   integer_options AS io 
                              WHERE  io.id = 6);

COMMIT;
`)
	if err != nil {
		return 0, err
	}

	ignited, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

  return extinguished + ignited, err 
}