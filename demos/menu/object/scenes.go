package object

import (
  "../../../backend"
)

func QuerySceneId(handle *backend.Handle, sceneName string) (int64, error) {
	var (
		err      error 
		row     *backend.Row 
		sceneId  int64 = -1
	)

  row, err = handle.QueryRow(`
SELECT s.id 
FROM   scenes AS s 
WHERE  s.name = ?;
`, sceneName)
	if err != nil {
		return sceneId, err
	}

	if row != nil {
  	err = row.Scan(&sceneId)  
	}

  return sceneId, err
}

func ScrollScene(handle *backend.Handle, sceneId int64, scrollX int32, scrollY int32) error {
	var err error 

	_, err = handle.Exec(`
UPDATE scenes 
SET    scroll_x = scroll_x + ? * scroll_speed,
       scroll_y = scroll_y + ? * scroll_speed 
WHERE  id = ?
	`, scrollX, scrollY, sceneId)

	return err
}