package main

import (
  "os"
  "fmt"
	"../.."
  "../../sdlex"
  "../../backend"
  "../../event"
  "./state/draw"
  "./state/scene"
  "./state/button"
  "./wrap"
)

func makeMenuProcessor(backendHandle *backend.Handle, sdlWrap *sdlex.SdlWrap, sdlexWrap sdlex.Wrap) (*event.Processor, error) {
  var (
    err             error
    menuWrap       *wrap.MenuWrap   = sdlexWrap.(*wrap.MenuWrap)
    eventProcessor *event.Processor
    init            scene.Init
    idle            button.Idle
  )

  eventProcessor = event.NewProcessor(sdlWrap)

  eventProcessor.AddProcess(event.NewProcess(draw.MakeIdle(backendHandle)))

  init, err = scene.MakeInit(sdlWrap.Handle(), menuWrap.Scene(), sdlWrap.Renderer())
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to make initializing scene state: %s\n", err)
    return eventProcessor, err
  }
  eventProcessor.AddProcess(event.NewProcess(init))

  idle, err = button.MakeIdle(1, backendHandle)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to make initial idle state: %s\n", err) 
  }
  eventProcessor.AddProcess(event.NewProcess(idle))

  idle, err = button.MakeIdle(2, backendHandle)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Failed to make initial idle state: %s\n", err) 
  }
  eventProcessor.AddProcess(event.NewProcess(idle))

  return eventProcessor, err
}

func menuSave(handle *backend.Handle) error {
  var err error 

  _, err = handle.Exec(`
INSERT OR ROLLBACK INTO save.integer_options  SELECT * FROM main.integer_options;
INSERT OR ROLLBACK INTO save.text_options     SELECT * FROM main.text_options;
INSERT OR ROLLBACK INTO save.boolean_options  SELECT * FROM main.boolean_options;
INSERT OR ROLLBACK INTO save.real_options     SELECT * FROM main.real_options;
INSERT OR ROLLBACK INTO save.objects          SELECT * FROM main.objects;
INSERT OR ROLLBACK INTO save.dots             SELECT * FROM main.dots;
INSERT OR ROLLBACK INTO save.lines            SELECT * FROM main.lines;
INSERT OR ROLLBACK INTO save.colors           SELECT * FROM main.colors;
INSERT OR ROLLBACK INTO save.states           SELECT * FROM main.states;
INSERT OR ROLLBACK INTO save.entities         SELECT * FROM main.entities;
INSERT OR ROLLBACK INTO save.scenes           SELECT * FROM main.scenes;
INSERT OR ROLLBACK INTO save.images           SELECT * FROM main.images;
INSERT OR ROLLBACK INTO save.sprites          SELECT * FROM main.sprites;
INSERT OR ROLLBACK INTO save.states_sprites   SELECT * FROM main.states_sprites;
INSERT OR ROLLBACK INTO save.entities_scenes  SELECT * FROM main.entities_scenes;
INSERT OR ROLLBACK INTO save.hitboxes         SELECT * FROM main.hitboxes;
`)

  return err
}

func menuLoad(handle *backend.Handle) error {
  var err error 

  _, err = handle.Exec(`
INSERT OR ROLLBACK INTO main.integer_options  SELECT * FROM save.integer_options;
INSERT OR ROLLBACK INTO main.text_options     SELECT * FROM save.text_options;
INSERT OR ROLLBACK INTO main.boolean_options  SELECT * FROM save.boolean_options;
INSERT OR ROLLBACK INTO main.real_options     SELECT * FROM save.real_options;
INSERT OR ROLLBACK INTO main.objects          SELECT * FROM save.objects;
INSERT OR ROLLBACK INTO main.dots             SELECT * FROM save.dots;
INSERT OR ROLLBACK INTO main.lines            SELECT * FROM save.lines;
INSERT OR ROLLBACK INTO main.colors           SELECT * FROM save.colors;
INSERT OR ROLLBACK INTO main.states           SELECT * FROM save.states;
INSERT OR ROLLBACK INTO main.entities         SELECT * FROM save.entities;
INSERT OR ROLLBACK INTO main.scenes           SELECT * FROM save.scenes;
INSERT OR ROLLBACK INTO main.images           SELECT * FROM save.images;
INSERT OR ROLLBACK INTO main.sprites          SELECT * FROM save.sprites;
INSERT OR ROLLBACK INTO main.states_sprites   SELECT * FROM save.states_sprites;
INSERT OR ROLLBACK INTO main.entities_scenes  SELECT * FROM save.entities_scenes;
INSERT OR ROLLBACK INTO main.hitboxes         SELECT * FROM save.hitboxes;
`)

  return err 
}

func main() {
  assemble.Run(
    menuSave, 
    menuLoad, 
    makeMenuProcessor, 
    wrap.MakeMenuWrap())
}