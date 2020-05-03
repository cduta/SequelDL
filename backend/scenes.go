package backend

func (handle *Handle) QuerySceneId(sceneName string) (int64, error) {
	var (
		err      error 
		row     *Row 
		sceneId  int64 = -1
	)

  row, err = handle.queryRow(`
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

func (handle *Handle) ScrollScene(sceneId int64, scrollX int32, scrollY int32) error {
	var err error 

	_, err = handle.exec(`
UPDATE scenes 
SET    scene_x = scene_x + ? * scroll_speed,
       scene_y = scene_y + ? * scroll_speed 
WHERE  id = ?
	`, scrollX, scrollY, sceneId)

	return err
}