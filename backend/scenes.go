package backend

import (
	"database/sql"
)

type Scenes struct {
	*Objects
}

func (handle *Handle) QuerySceneId(sceneName string) (int64, error) {
	var (
		err      error 
		row     *sql.Row 
		sceneId  int64
	)

  row, err = handle.queryRow(`
SELECT s.id 
FROM   scenes AS s 
WHERE  s.name = ?;
`, sceneName)
	if err != nil {
		return sceneId, err
	}

  err = row.Scan(&sceneId)  

  return sceneId, err
}

func (scenes Scenes) Close() {
  scenes.Objects.Close()
}