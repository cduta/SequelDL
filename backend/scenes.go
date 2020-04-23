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