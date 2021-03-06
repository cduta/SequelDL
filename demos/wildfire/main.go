package main

import (
	"../.."
  "../../sdlex"
  "../../backend"
  "../../event"
  "./wrap"
  "./state/idle"
  "./state/burn"
)

func makeWildfireProcessor(handle *backend.Handle, sdlWrap *sdlex.SdlWrap, sdlexWrap sdlex.Wrap) (*event.Processor, error) {
  var (
    err             error
    wildfireWrap   *wrap.WildfireWrap = sdlexWrap.(*wrap.WildfireWrap)
    eventProcessor *event.Processor
    idleState       idle.Idle 
    burnState       burn.Burn 
  )
  
  idleState, err = idle.MakeIdle(wildfireWrap.Particles(), handle)
  if err != nil {
    return eventProcessor, err
  }

  burnState, err = burn.MakeIdle(wildfireWrap.Particles(), handle)
  if err != nil {
    return eventProcessor, err
  }

  eventProcessor = event.NewProcessor(sdlWrap)
  eventProcessor.AddProcess(event.NewProcess(idleState))
  eventProcessor.AddProcess(event.NewProcess(burnState))

  return eventProcessor, err
}

func wildfireSave(handle *backend.Handle) error {
  var err error 

  _, err = handle.Exec(`
INSERT OR ROLLBACK INTO save.integer_options     SELECT * FROM main.integer_options;
INSERT OR ROLLBACK INTO save.text_options        SELECT * FROM main.text_options;
INSERT OR ROLLBACK INTO save.boolean_options     SELECT * FROM main.boolean_options;
INSERT OR ROLLBACK INTO save.real_options        SELECT * FROM main.real_options;
INSERT OR ROLLBACK INTO save.colors              SELECT * FROM main.colors;
INSERT OR ROLLBACK INTO save.color_ranges        SELECT * FROM main.color_ranges;
INSERT OR ROLLBACK INTO save.particles           SELECT * FROM main.particles;
INSERT OR ROLLBACK INTO save.states              SELECT * FROM main.states;
INSERT OR ROLLBACK INTO save.entities            SELECT * FROM main.entities;
INSERT OR ROLLBACK INTO save.entities_states     SELECT * FROM main.entities_states;
INSERT OR ROLLBACK INTO save.states_particles    SELECT * FROM main.states_particles;
INSERT OR ROLLBACK INTO save.old_entities_states SELECT * FROM main.old_entities_states;
`)

  return err
}

func wildfireLoad(handle *backend.Handle) error {
  var err error 

  _, err = handle.Exec(`
INSERT OR ROLLBACK INTO main.integer_options     SELECT * FROM save.integer_options;
INSERT OR ROLLBACK INTO main.text_options        SELECT * FROM save.text_options;
INSERT OR ROLLBACK INTO main.boolean_options     SELECT * FROM save.boolean_options;
INSERT OR ROLLBACK INTO main.real_options        SELECT * FROM save.real_options;
INSERT OR ROLLBACK INTO main.colors              SELECT * FROM save.colors;
INSERT OR ROLLBACK INTO main.color_ranges        SELECT * FROM save.color_ranges;
INSERT OR ROLLBACK INTO main.particles           SELECT * FROM save.particles;
INSERT OR ROLLBACK INTO main.states              SELECT * FROM save.states;
INSERT OR ROLLBACK INTO main.entities            SELECT * FROM save.entities;
INSERT OR ROLLBACK INTO main.entities_states     SELECT * FROM save.entities_states;
INSERT OR ROLLBACK INTO main.states_particles    SELECT * FROM save.states_particles;
INSERT OR ROLLBACK INTO main.old_entities_states SELECT * FROM save.old_entities_states;
`)

  return err 
}

func wildfireRendering(sdlWrap *sdlex.Wrap) error {
  return nil
}

func main() {
  assemble.Run(
    wildfireSave, 
    wildfireLoad, 
    makeWildfireProcessor, 
    wrap.MakeWildfireWrap())
}