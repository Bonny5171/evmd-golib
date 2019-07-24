package modelbase

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type MapJSON map[string]interface{}

func (j MapJSON) Value() (driver.Value, error) {
	p, err := json.Marshal(j)
	return p, err
}

func (j *MapJSON) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(v, &i)
	if err != nil {
		return err
	}

	*j, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}

	return nil
}
