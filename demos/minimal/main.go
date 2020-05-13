package main

import (
	"../.."
  "../../sdlex"
  "../../backend"
  "../../event"
  "./state/idle"
)

func makeMinimalProcessor(_ *backend.Handle, sdlWrap *sdlex.Wrap) (*event.Processor, error) {
  var (
    err             error
    eventProcessor *event.Processor
  )
  
  eventProcessor = event.NewProcessor(sdlWrap)
  eventProcessor.AddProcess(event.NewProcess(idle.MakeIdle()))

  return eventProcessor, err
}

func minimalSave(handle *backend.Handle) error {
  var err error 

  _, err = handle.Exec(`
INSERT OR ROLLBACK INTO save.integer_options  SELECT * FROM main.integer_options;
INSERT OR ROLLBACK INTO save.text_options     SELECT * FROM main.text_options;
INSERT OR ROLLBACK INTO save.boolean_options  SELECT * FROM main.boolean_options;
INSERT OR ROLLBACK INTO save.real_options     SELECT * FROM main.real_options;
`)

  return err
}

func minimalLoad(handle *backend.Handle) error {
  var err error 

  _, err = handle.Exec(`
INSERT OR ROLLBACK INTO main.integer_options  SELECT * FROM save.integer_options;
INSERT OR ROLLBACK INTO main.text_options     SELECT * FROM save.text_options;
INSERT OR ROLLBACK INTO main.boolean_options  SELECT * FROM save.boolean_options;
INSERT OR ROLLBACK INTO main.real_options     SELECT * FROM save.real_options;
`)

  return err 
}

func minimalRendering(sdlWrap *sdlex.Wrap) error {
  return nil
}

func main() {
  assemble.Run(minimalSave, minimalLoad, makeMinimalProcessor, minimalRendering)
}