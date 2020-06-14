package object

import (  
  "../../../backend"
)

type Particles struct {
  *backend.Rows
}

type Particle struct {
  backend.Position
  From        backend.Color 
  To          backend.Color
  RedrawDelay int32
  Id          int64
}

func QueryVisibleParticles(handle *backend.Handle) (*Particles, error) {
  var (
    err   error 
    rows *backend.Rows
  )

  rows, err = handle.QueryRows(`
SELECT p.id, 
       e.x + p.relative_x, e.y + p.relative_y, 
       c_from.r, c_from.g, c_from.b, c_from.a,
       c_to.r  , c_to.g  , c_to.b  , c_to.a, 
       cr.redraw_delay
FROM   entities         AS e,
       entities_states  AS es,
       states_particles AS sp,
       particles        AS p,
       color_ranges     AS cr,
       colors           AS c_from,
       colors           AS c_to
WHERE  e.visible
AND    e.id             = es.entity_id
AND    es.state_id      = sp.state_id
AND    sp.particle_id   = p.id
AND    p.color_range_id = cr.id
AND    cr.color_from    = c_from.id
AND    cr.color_to      = c_to.id
ORDER BY e.level, p.level;`)
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
    err         error
    position    backend.Position = backend.Position{}
    colorFrom   backend.Color    = backend.Color{}
    colorTo     backend.Color    = backend.Color{}
    redrawDelay int64
    particleId  int64 
  )

  if !particles.Rows.Next() {
    return nil, err
  }

  err = particles.Rows.Scan(
    &particleId, 
    &position.X , &position.Y, 
    &colorFrom.R, &colorFrom.G, &colorFrom.B, &colorFrom.A,
    &colorTo.R  , &colorTo.G  , &colorTo.B  , &colorTo.A  ,     
    &redrawDelay)

  return &Particle{ Position: position, From: colorFrom, To: colorTo, RedrawDelay: int32(redrawDelay), Id: particleId }, err
}
