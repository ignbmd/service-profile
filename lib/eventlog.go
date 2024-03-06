package lib

import (
	"log"

	"github.com/bytedance/sonic"
	"github.com/pandeptwidyaop/golog"
	"smartbtw.com/services/profile/db"
)

type EventLogBody struct {
	Version int            `json:"version"`
	Data    CreateEventLog `json:"data"`
}
type CreateEventLog struct {
	LogTitle   string `json:"log_title"`
	LogTags    string `json:"log_tags"`
	LogPayload string `json:"log_payload"`
	LogSource  string `json:"log_source"`
	LogDesc    string `json:"log_desc"`
	LogAction  string `json:"log_action"`
	LogType    string `json:"log_type"`
}

func LogEvent(title string, body []byte, action string, desc string, logtype string, tags string) bool {
	sd := EventLogBody{
		Version: 1,
		Data: CreateEventLog{
			LogTitle:   title,
			LogPayload: string(body),
			LogSource:  "Profile Service",
			LogTags:    tags,
			LogDesc:    desc,
			LogAction:  action,
			LogType:    logtype,
		},
	}

	jsonBody, err := sonic.Marshal(sd)
	if err != nil {
		text := "[Profile-EventLog] Failed to marshal summary to json"
		golog.Slack.ErrorWithData(text, body, err)
		return false
	}
	if db.Broker != nil {
		err = db.Broker.Publish("event-logging.created", "application/json", jsonBody)
		if err != nil {
			text := "[Profile-EventLog] Failed to publish event log"
			golog.Slack.ErrorWithData(text, body, err)
			log.Println(err)
			return false
		}
	} else {
		text := "[Profile-EventLog] Gorabbit is not initialized"
		golog.Slack.ErrorWithData(text, body, err)
		log.Println(err)
	}
	return true
}
