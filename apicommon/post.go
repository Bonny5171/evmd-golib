package apicommon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

//SendJobPost func
func SendJobPost(stack, orgID string, params map[string]string) error {
	jobAPIEndpoint := os.Getenv("JOB_API_ENDPOINT")
	if len(jobAPIEndpoint) < 1 {
		return fmt.Errorf("JOB_API_ENDPOINT not defined")
	}
	jobAPIToken := os.Getenv("JOB_API_TOKEN")
	if len(jobAPIToken) < 1 {
		return fmt.Errorf("JOB_API_TOKEN not defined")
	}
	client := &http.Client{}
	// body := map[string]string{"ETL_TABLE_NAME": etlTables}
	payload, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/org/%s/job/name/job_etl/push", jobAPIEndpoint, stack, orgID), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jobAPIToken))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errMessage interface{}
		errBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(errBody, &errMessage)
		if err != nil {
			return err
		}
		err = fmt.Errorf("%+v", errMessage)
		return err
	}
	return nil
}
