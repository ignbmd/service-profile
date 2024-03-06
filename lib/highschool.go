package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"smartbtw.com/services/profile/mockstruct"
)

func GetHighschoolStudent(schoolId string) (*mockstruct.SchoolLocation, error) {

	conn := os.Getenv("SERVICE_HIGHSCHOOL_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/school/%s", schoolId), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to highschool " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to highschool " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of highschool " + err.Error())
	}

	st := mockstruct.SchoolLocationRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of highschool " + errs.Error())
	}

	return &st.Data, nil
}
