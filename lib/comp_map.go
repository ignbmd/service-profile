package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/mockstruct"
)

func GetStudentProfileFromCompMap(smartbtwid uint) (*mockstruct.StudentProfileCompMapBody, error) {

	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/student-data/%d", smartbtwid), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map student profile " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map student profile " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map student profile " + err.Error())
	}

	st := mockstruct.StudentProfileCompMapRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map student profile " + errs.Error())
	}

	return &st.Data, nil
}

func GetSKDRankFromCompMap(majorId int, schoolId int, locationId *int) (*mockstruct.SKDRankCompMapBody, error) {

	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	url := fmt.Sprintf("/skd-rank?study_program_id=%d&school_id=%d", majorId, schoolId)
	if locationId != nil {
		url = fmt.Sprintf("/skd-rank?study_program_id=%d&school_id=%d&location_id=%d", majorId, schoolId, *locationId)
	}
	request, err = http.NewRequest("GET", conn+url, nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map student profile " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map student profile " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map student profile " + err.Error())
	}

	st := mockstruct.SKDRankCompMapRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map student profile " + errs.Error())
	}

	return &st.Data, nil
}

func GetCompetitionFromCompMap(compId uint) (*mockstruct.CompMapCompetition, error) {

	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/competition/%d", compId), nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map student profile " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map student profile " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map student profile " + err.Error())
	}

	st := mockstruct.CompetitionCompMapRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map student profile " + errs.Error())
	}

	return &st.Data, nil
}

func GetCompetitionMapData(majorId int, schoolId int, locationId *int, genders string) (*mockstruct.CompMapNewCompetition, error) {

	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	url := fmt.Sprintf("/competition/by-study-program/%d", majorId)
	if locationId != nil {
		if *locationId != 0 {
			url = fmt.Sprintf("/competition/by-study-program/%d/location/%d", majorId, *locationId)
		}
	}
	if schoolId == 5 || schoolId == 6 {
		url = fmt.Sprintf("%s?tags=%s", url, genders)
	}
	request, err = http.NewRequest("GET", conn+url, nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map student profile " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map student profile " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map student profile " + err.Error())
	}

	st := mockstruct.CompetitionNewCompMapRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map student profile " + errs.Error())
	}

	if len(st.Data) < 1 {
		return nil, errors.New("data not found")
	}

	return &st.Data[0], nil
}

func GetCompetitionDataPTK(majorId uint, locationId uint, gender string, polbitType string) (*mockstruct.CompetitionPTKPTNData, error) {

	res, err := _fetchCompetitionDataPTK(majorId, locationId, gender, polbitType)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	yearIndex := 0
	quotaYearIndex := 0

	resComp := mockstruct.CompetitionPTKPTNData{}

	if len(*res) < 1 {
		return &resComp, nil
	}

	for idx, k := range *res {
		if k.Quota > 0 {
			quotaYearIndex = idx
			break
		}
	}
	for idx, k := range *res {
		if k.Registered > 0 {
			yearIndex = idx
			break
		}
	}

	if majorId == 348 || majorId == 349 {
		if (*res)[quotaYearIndex].ID > 4266 {
			resComp.CompetitionType = &(*res)[quotaYearIndex].CompetitionType
		}
	}

	resComp.MajorQuota = uint32((*res)[quotaYearIndex].Quota)
	resComp.MajorRegistered = uint32((*res)[yearIndex].Registered)
	resComp.MajorYear = uint32((*res)[yearIndex].Year)
	resComp.MajorQuotaYear = uint32((*res)[quotaYearIndex].Year)
	resComp.MajorOldQuota = uint32((*res)[yearIndex].Quota)

	return &resComp, nil
}

func _fetchCompetitionDataPTK(majorId uint, locationId uint, gender string, polbitType string) (*[]mockstruct.CompetitionPTK, error) {
	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}

	url := fmt.Sprintf("%s/competition/by-study-program/%d", conn, majorId)

	if majorId != 348 && majorId != 349 {
		gender = ""
	}

	isAffirmasi := false

	if strings.Contains(strings.ToUpper(polbitType), "_AFIRMASI") {
		isAffirmasi = true
	}

	if locationId != 0 && strings.Contains(strings.ToUpper(polbitType), "DAERAH") {
		url = fmt.Sprintf("%s/location/%d?tags=%s&is_afirm=%v", url, locationId, gender, isAffirmasi)
	} else {
		url = fmt.Sprintf("%s/%d?tags=%s", url, locationId, gender)
	}

	request, err := http.NewRequest("GET", url, nil)
	// TODO: check err
	if err != nil {
		golog.Slack.Error("creating request to comp map", err)
		return nil, errors.New("server error")
	}

	request.Close = true

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		golog.Slack.Error("doing request to comp map", err)
		return nil, errors.New("server error")
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		golog.Slack.Error("reading response of comp map", err)
		return nil, errors.New("server error")
	}

	if resp.StatusCode >= 300 {
		type HttpResponseBody struct {
			Message any `json:"message"`
		}

		var resp HttpResponseBody

		_ = sonic.Unmarshal(b, &resp)

		return nil, fmt.Errorf("%v", resp.Message)
	}

	type responseBody struct {
		Message any `json:"message"`
		Data    *[]mockstruct.CompetitionPTK
	}

	var decoded responseBody

	errs := sonic.Unmarshal(b, &decoded)
	if errs != nil {
		return nil, errors.New("unmarshalling response of location " + errs.Error())
	}

	return decoded.Data, nil
	// TODO: check err
}

func GetCompetitonPTN(studyProgramID uint) (mockstruct.CompetitionPTN, error) {
	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return mockstruct.CompetitionPTN{}, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+fmt.Sprintf("/passing-grade/by-study-program/%d", studyProgramID), nil)

	// TODO: check err
	if err != nil {
		return mockstruct.CompetitionPTN{}, errors.New("creating request to comp map student profile " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return mockstruct.CompetitionPTN{}, errors.New("doing request to comp map student profile " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return mockstruct.CompetitionPTN{}, errors.New("reading response of comp map student profile " + err.Error())
	}

	st := mockstruct.CompetitionPTNRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return mockstruct.CompetitionPTN{}, errors.New("unmarshalling response of comp map student profile " + errs.Error())
	}

	return st.Data, nil
}

func fetchCompetitionFormationCPNS(req mockstruct.GetCompetitionCPNS) ([]mockstruct.CompetitionFormationCPNS, error) {
	conn := os.Getenv("SERVICE_COMP_MAP_V2_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	bd := map[string]any{
		"formation_type": req.FormationType,
		"position_id":    req.PositionID,
		"formation_code": req.FormationCode,
	}
	ns, _ := sonic.Marshal(bd)

	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("POST", conn+"/competition-cpns-chance", bytes.NewBuffer(ns))

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to comp map student profile " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to comp map student profile " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of comp map student profile " + err.Error())
	}
	st := mockstruct.CompetitionFormationCPNSRequest{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of comp map student profile " + errs.Error())
	}

	return st.Data, nil
}

func GetCompetitionFormationCPNS(req mockstruct.GetCompetitionCPNS) (*mockstruct.CompetitionFormationCPNS, error) {
	chance, err := fetchCompetitionFormationCPNS(req)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(chance))

	valMajorQuota := 0
	valMajorRegistrant := 0
	valMajorQuotaYear := 0
	if len(chance) > 0 {
		valMajorQuota = int(chance[0].Quota)
		valMajorRegistrant = int(chance[0].Registered)
		valMajorQuotaYear = int(chance[0].Year)
		// valMajorCompYear = int(chance[0].Year)

		for _, k := range chance {
			if valMajorRegistrant > 0 {
				break
			}
			if k.Registered > 0 {
				valMajorRegistrant = int(k.Registered)
				// valMajorCompYear = int(k.Year)
			}
		}
		// thg := math.Round(float64(valMajorRegistrant) / float64(valMajorQuota))
		// valMajorChances = fmt.Sprintf("1:%.0f", thg)
	}

	res := mockstruct.CompetitionFormationCPNS{
		PositionID: req.PositionID,
		Quota:      uint(valMajorQuota),
		Registered: uint(valMajorRegistrant),
		Year:       uint(valMajorQuotaYear),
	}
	return &res, nil
}
