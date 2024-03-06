package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
)

func GetStudentPresence(smartbtwID uint) (map[string]any, error) {

	conn := os.Getenv("SERVICE_LEARNING_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/student-presence/student-summary/for-report/%d", smartbtwID), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to exam " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to exam " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of exam " + err.Error())
	}

	type GetStudentPresenceResponse struct {
		Data map[string]any `json:"data"`
	}
	st := GetStudentPresenceResponse{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of exam " + errs.Error())
	}

	return st.Data, nil
}
