package backend

import (
	"fmt"
)

type Scenes struct {
	*Objects
}

func (handle Handle) QuerySceneId(sceneName string) (int64, error) {
	var (
		err      error 
		objects *Objects 
		scenes   Scenes
		sceneId  int64
	)

  objects, err = handle.queryObjects(`
SELECT s.id 
FROM   scenes AS s 
WHERE  s.name = ?;
`, sceneName)
	if err != nil {
		return sceneId, err
	}

	scenes = Scenes{ Objects: objects }
  defer scenes.Close()

  if !scenes.Objects.rows.Next() {
    return sceneId, fmt.Errorf("Could not find the id of a sprite with the name: %s", sceneName) 
  }

  err = scenes.Objects.rows.Scan(&sceneId)  

  return sceneId, err
}

func (scenes Scenes) Close() {
  scenes.Objects.Close()
}