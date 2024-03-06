package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
)

type User struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
type GetStudentPresenceResponse struct {
	Users User `json:"users"`
}

func GetUserSSO(userID string) (*User, error) {

	conn := os.Getenv("SERVICE_SSO_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)

	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/single-user/%s", userID), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-SSO-Token", os.Getenv("TOKEN_SSO_V2"))

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

	st := GetStudentPresenceResponse{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of exam " + errs.Error())
	}

	return &st.Users, nil
}
