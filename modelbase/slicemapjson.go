package modelbase

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SliceMapJSON []map[string]interface{}

func (j SliceMapJSON) Value() (driver.Value, error) {
	p, err := json.Marshal(j)
	return p, err
}

func (j *SliceMapJSON) Scan(value interface{}) error {
	v, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	// var i interface{}
	err := json.Unmarshal(v, &j)
	if err != nil {
		return err
	}

	return nil
}
