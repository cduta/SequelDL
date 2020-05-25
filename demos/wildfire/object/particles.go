package object

import (  
  "../../../backend"
)

type Particles struct {
  *backend.Rows
}

type Particle struct {
  backend.Position
  backend.Color 
  Id int64
}

func QueryFireParticles(handle *backend.Handle) (*Particles, error) {
  var (
    err   error 
    rows *backend.Rows
  )

  rows, err = handle.QueryRows(`
SELECT p.id, 
		   e.x + p.relative_x + offset_x.column1, 
		   e.y + p.relative_y + offset_y.column1, 
		   c_from.r + abs(random() % (abs(c_to.r-c_from.r)+1)), 
		   c_from.g + abs(random() % (abs(c_to.g-c_from.g)+1)), 
		   c_from.b + abs(random() % (abs(c_to.b-c_from.b)+1)), 
		   c_from.a + abs(random() % (abs(c_to.a-c_from.a)+1))
FROM   entities         AS e,
       entities_states  AS es,
       states_particles AS sp,
       particles        AS p,
       color_ranges     AS cr,
       colors           AS c_from,
       colors           AS c_to, ( 
         VALUES (-1),( 0),( 1)
       ) AS offset_x, (
         VALUES (-1),( 0),( 1)
       ) AS offset_y
WHERE  e.visible
AND    sp.state_id      = 2
AND    e.id             = es.entity_id
AND    es.state_id      = sp.state_id
AND    sp.particle_id   = p.id
AND    p.color_range_id = cr.id
AND    cr.color_from    = c_from.id
AND    cr.color_to      = c_to.id
AND    (offset_x.column1 = 0 OR offset_y.column1 = 0);`)
  if rows == nil || err != nil {
    return nil, err 
  }

  return &Particles{ Rows: rows }, err
}

func (particles Particles) Close() {
  particles.Rows.Close()
}

func (particles Particles) Next() (*Particle, error) {
  var (
    err        error
    position   backend.Position = backend.Position{}
    color      backend.Color    = backend.Color{}
    particleId int64 
  )

  if !particles.Rows.Next() {
    return nil, err
  }

  err = particles.Rows.Scan(&particleId, &position.X, &position.Y, &color.R, &color.G, &color.B, &color.A)

  return &Particle{ Position: position, Color: color, Id: particleId}, err
}
