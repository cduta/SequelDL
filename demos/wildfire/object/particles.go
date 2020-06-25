package object

import (  
  "../../../backend"
)

type Particles struct {
  *backend.Rows
}

type Particle struct {
  backend.Position
  From         backend.Color 
  To           backend.Color
  CanBeRedrawn bool 
  RedrawIn     uint32
  RedrawDelay  uint32
  Id           int64
}

func (particle *Particle) AdvanceRedrawDelay() bool {
  var redraw bool = false

  if !particle.CanBeRedrawn {
    return redraw
  }

  particle.RedrawIn--

  if particle.RedrawIn == 0 {
    redraw = true
    particle.RedrawIn = particle.RedrawDelay
  }

  return redraw
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

func QueryOldParticles(handle *backend.Handle) (*Particles, error) {
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
FROM   (SELECT oes.entity_id AS entity_id, 
               oes.state_id  AS state_id
        FROM   old_entities_states AS oes) AS oes,
       entities            AS e,
       states_particles    AS sp,
       particles           AS p,
       color_ranges        AS cr,
       colors              AS c_from,
       colors              AS c_to
WHERE  e.visible
AND    e.id             = oes.entity_id
AND    oes.state_id     = sp.state_id
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

func QueryChangedParticles(handle *backend.Handle) (*Particles, error) {
  var (
    err   error 
    rows *backend.Rows
  )

  rows, err = handle.QueryRows(`
SELECT p.id, 
       e2.x + p.relative_x, e2.y + p.relative_y, 
       c_from.r, c_from.g, c_from.b, c_from.a,
       c_to.r  , c_to.g  , c_to.b  , c_to.a, 
       cr.redraw_delay
FROM   (SELECT DISTINCT e.x + p.relative_x AS x, e.y + p.relative_y AS y
        FROM   entities            AS e,
               old_entities_states AS oes,
               states_particles    AS sp,
               particles           AS p
        WHERE  e.id           = oes.entity_id 
        AND    oes.state_id   = sp.state_id
        AND    sp.particle_id = p.id
          UNION 
        SELECT DISTINCT e.x + p.relative_x AS x, e.y + p.relative_y AS y
        FROM   old_entities_states AS oes,
               entities            AS e,
               entities_states     AS es,
               states_particles    AS sp,
               particles           AS p
        WHERE  e.id           = oes.entity_id 
        AND    es.entity_id   = e.id
        AND    es.state_id    = sp.state_id
        AND    sp.particle_id = p.id) AS e1,
       entities            AS e2,
       entities_states     AS es,
       states_particles    AS sp,
       particles           AS p,
       color_ranges        AS cr,
       colors              AS c_from,
       colors              AS c_to
WHERE  e2.visible 
AND    (e1.x, e1.y)     = (e2.x + p.relative_x, e2.y + p.relative_y)
AND    e2.id            = es.entity_id
AND    es.state_id      = sp.state_id
AND    sp.particle_id   = p.id 
AND    p.color_range_id = cr.id 
AND    cr.color_from    = c_from.id 
AND    cr.color_to      = c_to.id 
ORDER BY p.level, e2.level;`)
  if rows == nil || err != nil {
    return nil, err 
  }

  return &Particles{ Rows: rows }, err
}

func ClearOldStates(handle *backend.Handle) error {
  var err error

  _, err = handle.Exec(`
BEGIN IMMEDIATE;
DELETE FROM old_entities_states;
COMMIT;
`)

  return err
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

  return &Particle{ Position: position, From: colorFrom, To: colorTo, CanBeRedrawn: redrawDelay > 0, RedrawIn: uint32(redrawDelay), RedrawDelay: uint32(redrawDelay), Id: particleId }, err
}
