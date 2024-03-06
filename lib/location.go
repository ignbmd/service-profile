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

func GetLocationByID(id uint) (*mockstruct.Location, error) {

	conn := os.Getenv("SERVICE_LOCATION_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/location/by-ids?ids=1&ids=%d", id), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to location " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to location " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of location " + err.Error())
	}

	st := mockstruct.LocationRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of location " + errs.Error())
	}

	if len(st.Data) > 0 {
		return &st.Data[0], nil
	}

	return nil, nil
}
