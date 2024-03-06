package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/mockstruct"
	requests "smartbtw.com/services/profile/request"
)

func GetStagesChallengeList(program string) ([]mockstruct.Stages, error) {
	refParent := "stages"
	allS := []mockstruct.Stages{}

	rs := map[string]mockstruct.Stages{}
	refAllSt := fmt.Sprintf("all-stages-%s", strings.ToLower(program))

	ref := db.FbDBMulti["STAGES"].NewRef(refParent)
	err := ref.Child(refAllSt).Get(db.Ctx, &rs)
	if err != nil {
		return []mockstruct.Stages{}, err
	}

	for _, k := range rs {
		if k.ModuleType != "PREMIUM_TRYOUT" {
			continue
		}
		allS = append(allS, k)
	}

	sort.Slice(allS,
		func(i, j int) bool {
			return allS[i].Stage < allS[j].Stage
		})

	return allS, nil
}

func GetAllStudentStageClass(typ, typStg string) (*requests.StudentStagesClassResult, error) {
	conn := os.Getenv("SERVICE_STAGES_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}
	var (
		request *http.Request
		err     error
	)
	request, err = http.NewRequest("GET", conn+"/stages-regular-class-all/"+typ+"/"+typStg, nil)

	// TODO: check err
	if err != nil {
		return nil, errors.New("creating request to product " + err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, errors.New("doing request to stages " + err.Error())
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("reading response of stages " + err.Error())
	}

	st := requests.StudentStagesClassResult{}

	errs := sonic.Unmarshal(b, &st)
	if errs != nil {
		return nil, errors.New("unmarshalling response of stages" + errs.Error())
	}

	return &st, nil
}
