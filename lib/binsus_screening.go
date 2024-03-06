package lib

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/mockstruct"
)

func FetchStudentBinsusSummary(smId int, year int) (*mockstruct.BinsusScreeningSummary, error) {
	conn := os.Getenv("SERVICE_BINSUS_SCREENING_HOST")
	if conn == "" {
		return nil, errors.New("host is null")
	}

	url := fmt.Sprintf("%s/summary/by-student-ids/%d/%d", conn, smId, year)

	request, err := http.NewRequest("GET", url, nil)
	// TODO: check err
	if err != nil {
		golog.Slack.Error("creating request to binsus screening", err)
		return nil, errors.New("server error")
	}

	request.Close = true

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		golog.Slack.Error("doing request to binsus screening", err)
		return nil, errors.New("server error")
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		golog.Slack.Error("reading response of binsus screening", err)
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
		Data    *mockstruct.BinsusScreeningSummary
	}

	var decoded responseBody

	errs := sonic.Unmarshal(b, &decoded)
	if errs != nil {
		return nil, errors.New("unmarshalling response of location " + errs.Error())
	}

	return decoded.Data, nil
	// TODO: check err
}
