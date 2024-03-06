package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/joho/godotenv"
	"github.com/pandeptwidyaop/golog"
	"github.com/pandeptwidyaop/gorabbit"
	"github.com/thatisuday/commando"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/lib/scripts"
	"smartbtw.com/services/profile/request"
)

func main() {
	commando.
		SetExecutableName("profile-cli").
		SetVersion("0.0.1").
		SetDescription("This CLI tool for smartbtw profile services")

	commando.
		Register("historyptk:sync").
		SetDescription("sync start end date for history ptk data").
		SetShortDescription("for sync data history ptk").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			res, err := lib.GetStudentHistoryPTKOnlyStage()
			if err != nil {
				fmt.Println(err)
			}

			for _, v := range res {
				bdy := request.HistoryPTKSendToExamResult{
					SmartbtwID: v.SmartBtwID,
					TaskID:     v.TaskID,
					Repeat:     v.Repeat,
					Program:    "skd",
				}
				msgBdy := request.MessageHistoryPTKSendToExamResult{
					Version: 1,
					Data:    bdy,
				}
				msgJson, err := sonic.Marshal(msgBdy)
				if err != nil {
					fmt.Println("error marshal message body")
				}

				if err = db.Broker.Publish(
					"result.sync-profile",
					"application/json",
					[]byte(msgJson), // message to publish
				); err != nil {
					fmt.Println("error on publishing mq for sync data" + err.Error())
				}

				fmt.Printf("Message sent to RabbitMQ: %d\n", v.SmartBtwID)
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("historyptn:sync").
		SetDescription("sync start end date for history ptn data").
		SetShortDescription("for sync data history ptn").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			res, err := lib.GetStudentHistoryPTNOnlyStage()
			if err != nil {
				fmt.Println(err)
			}

			for _, v := range res {
				bdy := request.HistoryPTKSendToExamResult{
					SmartbtwID: v.SmartBtwID,
					TaskID:     v.TaskID,
					Repeat:     v.Repeat,
					Program:    "utbk_irt",
				}
				msgBdy := request.MessageHistoryPTKSendToExamResult{
					Version: 1,
					Data:    bdy,
				}
				msgJson, err := sonic.Marshal(msgBdy)
				if err != nil {
					fmt.Println("error marshal message body")
				}

				if err = db.Broker.Publish(
					"result.sync-profile",
					"application/json",
					[]byte(msgJson), // message to publish
				); err != nil {
					fmt.Println("error on publishing mq for sync data" + err.Error())
				}

				fmt.Printf("Message sent to RabbitMQ: %d\n", v.SmartBtwID)
			}
			fmt.Printf("Command finished\n")
		})
	commando.
		Register("student-profile-ptk:sync").
		SetDescription("sync start end date for student profile PTK (PKN-STAN) data").
		SetShortDescription("for sync data student profile PTK (PKN-STAN)").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			err := lib.SyncSchoolPTKUnmatch()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Command finished\n")
		})

	commando.
		Register("classmember:sync").
		SetDescription("sync start end date for history ptn data").
		SetShortDescription("for sync data history ptn").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			res, err := lib.GetAllClassMember()
			if err != nil {
				fmt.Println(err)
			}

			for _, v := range res {
				err = lib.SyncClassMemberToElastic(uint(v.SmartbtwID))
				if err != nil {
					fmt.Println(err)
				}

				fmt.Printf("student id: %d and classroom_id : %s generated\n", v.SmartbtwID, v.ClassroomID.Hex())
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("historyptk:sync-timestamp").
		SetDescription("sync timestamp for history ptk data").
		SetShortDescription("for sync data history ptk").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {

			fmt.Printf("Obtaining history PTK data\n")

			res, err := lib.GetAllStudentHistoryPTK()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Obtained %d of history PTK data\n", len(res))

			for _, v := range res {
				fmt.Printf("Procesing SmartBTW ID: %d with TaskID : %d and Repeat : %d\n", v.SmartBtwID, v.TaskID, v.Repeat)
				bdy := request.BodyUpdateStudentDuration{
					SmartbtwID: v.SmartBtwID,
					TaskID:     v.TaskID,
					Repeat:     v.Repeat,
					Start:      v.CreatedAt,
					End:        v.UpdatedAt,
				}

				err := lib.UpdateTimestampHistoryPtkElastic(&bdy)
				if err != nil {
					fmt.Println("failed to sync", err)
				}
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("historyptn:sync-timestamp").
		SetDescription("sync timestamp for history ptn data").
		SetShortDescription("for sync data history ptn").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {

			fmt.Printf("Obtaining history PTN data\n")
			res, err := lib.GetAllStudentHistoryPTN()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Obtained %d of history PTN data\n", len(res))

			for _, v := range res {
				fmt.Printf("Procesing SmartBTW ID: %d with TaskID : %d and Repeat : %d\n", v.SmartBtwID, v.TaskID, v.Repeat)
				bdy := request.BodyUpdateStudentDuration{
					SmartbtwID: v.SmartBtwID,
					TaskID:     v.TaskID,
					Repeat:     v.Repeat,
					Start:      v.CreatedAt,
					End:        v.UpdatedAt,
				}
				err := lib.UpdateTimestampHistoryPtnElastic(&bdy)
				if err != nil {
					fmt.Println("failed to sync", err)
				}
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("profileptk:update-ipdn-scores").
		SetDescription("update ipdn ptk profile target scores").
		SetShortDescription("for sync data student profile ptk").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			err := scripts.UpdateIPDNHistoryScore()
			if err != nil {
				fmt.Println("failed to sync", err)
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("profileptk:update-std-scores").
		SetDescription("update std ptk profile target scores").
		SetShortDescription("for sync data student profile ptk").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			err := scripts.UpdateSTDHistoryScore()
			if err != nil {
				fmt.Println("failed to sync", err)
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("profileptk:resync-ptk").
		AddArgument("polbit_type", "the polbit type target to sync to", "").
		SetDescription("sync ptk profile target scores").
		SetShortDescription("for sync data student profile ptk").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			if m1["polbit_type"].Value == "" {
				fmt.Println("usage: [polbit_type]")
				return
			}
			err := scripts.ResyncPolbitHistoryScore(m1["polbit_type"].Value)
			if err != nil {
				fmt.Println("failed to sync", err)
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("profile:sync-branch-code").
		AddArgument("start_offset", "the start offset smartbtwid target to sync to", "").
		SetDescription("sync student profile branch codes to elastic").
		SetShortDescription("sync student profile branch codes to elastic").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			if m1["start_offset"].Value == "" {
				fmt.Println("usage: [start_offset]")
				return
			}
			offset, err := strconv.Atoi(m1["start_offset"].Value)
			if err != nil {
				fmt.Println("Error when insert json: ", err)
				return
			}
			err = scripts.SyncStudentProfileBranchCodes(offset)
			if err != nil {
				fmt.Println("failed to sync", err)
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("profile:sync-study-program-ssn").
		SetDescription("sync student profile study program SSN school").
		SetShortDescription("sync student profile study program SSN school").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			ids := []int{311, 312, 313}
			py := request.UpdateSpecificStudyProgram{
				MajorID:    ids,
				TargetType: "PTK",
				NewMajorID: 311,
				MajorName:  "D-IV Rekayasa Perangkat Keras Kripto, D-IV Rekayasa Keamanan Siber, D-IV Rekayasa Sistem Kriptografi",
			}
			err := lib.UpdateSpecificStudyProgram(context.Background(), &py)
			if err != nil {
				fmt.Println("failed to sync", err)
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("profileptk:resync-ptk-new").
		SetDescription("sync ptk profile target scores").
		SetShortDescription("for sync data student profile ptk").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			err := scripts.ResyncCentralTargetScore()
			if err != nil {
				fmt.Println("failed to sync", err)
			}
			fmt.Printf("Command finished\n")
		})

	commando.
		Register("history-score:resync-additional-data").
		AddFlag("program", "choose program to process", commando.String, "").
		AddFlag("force", "force resync", commando.Bool, false).
		SetDescription("sync history score additional data").
		SetShortDescription("sync history score additional data").
		SetAction(func(m1 map[string]commando.ArgValue, m2 map[string]commando.FlagValue) {
			if m2["program"].Value == "" {
				log.Fatal("program required")
			}

			err := scripts.ResyncStudentHistoryAdditionalData(strings.ToLower(fmt.Sprintf("%v", m2["program"].Value)), m2["force"].Value.(bool))
			if err != nil {
				fmt.Println("failed to sync", err)
			}
			fmt.Printf("Command finished\n")
		})
	commando.Parse(nil)
}

func init() {
	if godotenv.Load(".env") != nil {
		log.Println("Unable to load .env file, using global varibale")
	}
	connection := os.Getenv("MONGODB_CONNECTION")
	database := os.Getenv("MONGODB_DATABASE")
	db.Connect(connection, database)
	db.NewElastic()

	var err error
	conn := os.Getenv("RABBITMQ_CONNECTION")
	app := os.Getenv("APP_NAME")
	if conn == "" {
		golog.Slack.Error("Rabbit MQ URL Connection not set", nil)
		log.Panic("no varible found for RABBITMQ_CONNECTION")
	}

	if app == "" {
		log.Panic("app name not initialize")
	}

	db.Broker, err = gorabbit.New(conn, app, "GLOBAL_X")

	if err != nil {
		golog.Slack.Error("Error when create new connection to Rabbit MQ server", err)
		log.Panic(err)
	}

	db.Broker.Connect()

}
