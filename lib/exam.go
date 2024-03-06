package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/helpers"
	"smartbtw.com/services/profile/mockstruct"
)

func GetTryoutCodeBySchoolID(schoolID string) ([]mockstruct.GetUKACodeByschoolID, error) {

	conn := os.Getenv("SERVICE_NEW_EXAM_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/client/tryout/code/%s", schoolID), nil)

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

	st := mockstruct.GetUKACodeBySchoolIDResponse{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of exam " + errs.Error())
	}

	return st.Data, nil
}

func GetPackageByID(packageID uint) (*mockstruct.Packages, error) {

	conn := os.Getenv("SERVICE_NEW_EXAM_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/client/packages/id/%d", packageID), nil)

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

	st := mockstruct.GetPackageByIDResponse{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of exam " + errs.Error())
	}

	return &st.Data, nil
}

func GetPackageCPNSByID(packageID uint) (*mockstruct.Packages, error) {

	conn := os.Getenv("SERVICE_EXAM_CPNS_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/packages/caches/%d", packageID), nil)

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

	st := mockstruct.GetPackageByIDResponse{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of exam " + errs.Error())
	}

	return &st.Data, nil
}

func SendHistoryStage(smID uint, packageID uint, stgType string) error {
	var pkg *mockstruct.Packages
	if stgType == "CPNS" {
		pkgData, err := GetPackageCPNSByID(packageID)
		if err != nil {
			return err
		}
		pkg = pkgData
	} else {
		pkgData, err := GetPackageByID(packageID)
		if err != nil {
			return err
		}
		pkg = pkgData
	}
	if helpers.ArrayContainsCS(pkg.Tags, "STAGE_") {
		stagePayload := mockstruct.SendHistoryStageBody{
			Version: 2,
			Data: map[string]interface{}{
				"smartbtw_id": smID,
				"package_id":  packageID,
				"stage_type":  stgType,
			},
		}

		stageBody, err := sonic.Marshal(stagePayload)
		if err != nil {
			return errors.New("marshalling " + err.Error())
		}

		if err = db.Broker.Publish(
			"historystage.created",
			"application/json",
			[]byte(stageBody), // message to publish
		); err != nil {
			return errors.New("sending to rabbitmq " + err.Error())
		}
	}
	return nil
}
