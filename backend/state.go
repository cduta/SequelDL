package backend

func (handle *Handle) ChangeState(stateName string, entityId int64) error {
  var err error 

  _, err = handle.Exec(`
UPDATE entities   
SET    state_id = (SELECT st.id
                   FROM   states AS st 
                   WHERE  st.name = ?)
WHERE  id = ?;
`, stateName, entityId)

  return err
}