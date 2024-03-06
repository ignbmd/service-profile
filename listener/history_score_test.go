package listener_test

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"smartbtw.com/services/profile/lib"
	"smartbtw.com/services/profile/listener"
	"smartbtw.com/services/profile/models"
	"smartbtw.com/services/profile/request"
)

// func Test_UpdateScorePTK(t *testing.T) {
// 	Init()

// 	sPhoto := "http://www.google.com"
// 	sPayload := request.BodyMessageCreateUpdateStudent{
// 		ID:    200001,
// 		Name:  "John Doe 2",
// 		Email: "john_doe2@email.com",
// 		Photo: &sPhoto,
// 	}
// 	_, errS := lib.UpsertStudent(&sPayload)
// 	assert.Nil(t, errS)

// 	payload2 := request.CreateWallet{
// 		SmartbtwID: 200001,
// 		Point:      0,
// 		Type:       string(models.BONUS),
// 	}

// 	_, err2 := lib.CreateWallet(&payload2)
// 	if err2 != nil {
// 		assert.NotNil(t, err2)
// 	} else {
// 		assert.Nil(t, err2)
// 	}

// 	json := `{
// 		"version" : 2,
// 		"data": {
// 			"smartbtw_id": 200001,
// 			"school_id": 2,
// 			"major_id":2,
// 			"school_name": "UI",
// 			"major_name": "Kedokteran",
// 			"target_score": 400,
// 			"target_type": "PTK"
// 		}
// 	}`

// 	body := []byte(json)
// 	payload1, _ := request.UnmarshalMessageStudentTargetBody(body)
// 	_, err := lib.GetStudentTargetByCustom(payload1.Data.SmartbtwID, "PTK")
// 	if err != nil {
// 		_, err2 := lib.CreateStudentTarget(&payload1.Data)
// 		assert.NotNil(t, err)
// 		assert.Nil(t, err2)
// 	} else {
// 		assert.Nil(t, err)
// 	}

// 	payload := []byte(`
// 	{
// 		"version": 2,
// 		"data": {
// 			"smartbtw_id": 200001,
// 			"task_id": 1,
// 			"module_code": "MD-102",
// 			"module_type": "uka_premium",
// 			"tiu": 155,
// 			"tiu_pass_status": true,
// 			"tiu_passing_grade": 80,
// 			"tkp": 194,
// 			"tkp_pass_status": true,
// 			"tkp_passing_grade": 156,
// 			"total": 523,
// 			"twk": 174,
// 			"twk_pass_status": true,
// 			"is_all_passed": false,
// 			"twk_passing_grade": 65,
// 			"repeat": 1,
// 			"exam_name": "GANTENG 2",
// 			"start" : "2023-01-10T06:39:33.595+00:00",
// 			"end" : "2023-01-10T06:40:33.595+00:00",
// 			"is_live": true
// 		}
// 	}
// 	`)

// 	msg := amqp.Delivery{
// 		RoutingKey:   "history-ptk.created",
// 		Body:         payload,
// 		Acknowledger: &MockAcknowledger{},
// 	}
// 	fmt.Printf("%s", payload)
// 	assert.True(t, listener.ListenScoreHistoryBinding(&msg))
// }

// func Test_UpdateScorePTN(t *testing.T) {
// 	Init()

// 	sPhoto := "http://www.google.com"
// 	sPayload := request.BodyMessageCreateUpdateStudent{
// 		ID:    200008,
// 		Name:  "John Doe 7",
// 		Email: "john_doe5@email.com",
// 		Photo: &sPhoto,
// 	}
// 	_, errS := lib.UpsertStudent(&sPayload)
// 	assert.Nil(t, errS)

// 	payload2 := request.CreateWallet{
// 		SmartbtwID: 200008,
// 		Point:      0,
// 		Type:       string(models.BONUS),
// 	}

// 	_, err2 := lib.CreateWallet(&payload2)
// 	if err2 != nil {
// 		assert.NotNil(t, err2)
// 	} else {
// 		assert.Nil(t, err2)
// 	}

// 	json := `{
// 		"version" : 2,
// 		"data": {
// 			"smartbtw_id": 200008,
// 			"school_id": 1,
// 			"major_id":1,
// 			"school_name": "UI",
// 			"major_name": "Kedokteran",
// 			"target_score": 400,
// 			"target_type": "PTN"
// 		}
// 	}`

// 	body := []byte(json)
// 	payload1, _ := request.UnmarshalMessageStudentTargetBody(body)
// 	_, err := lib.GetStudentTargetByCustom(payload1.Data.SmartbtwID, "PTN")
// 	if err != nil {
// 		_, err2 := lib.CreateStudentTarget(&payload1.Data)
// 		assert.Nil(t, err2)
// 		assert.NotNil(t, err)
// 	} else {
// 		assert.Nil(t, err)
// 	}

// 	payload := []byte(`
// 	{
// 		"version": 1,
// 		"data": {
// 			"smartbtw_id": 200008,
// 			"task_id": 100,
// 			"module_code": "MD-808",
// 			"module_type": "uka_premium",
// 			"potensi_kognitif": 300,
// 			"penalaran_matematika": 300,
// 			"literasi_bahasa_indonesia": 300,
// 			"literasi_bahasa_inggris": 300,
// 			"total": 900,
// 			"repeat": 1,
// 			"exam_name": "TPS IRT 1",
// 			"start" : "2023-01-10T06:39:33.595+00:00",
// 			"end" : "2023-01-10T06:40:33.595+00:00"
// 		}
// 	}
// 	`)

// 	msg := amqp.Delivery{
// 		RoutingKey:   "history-ptn.created",
// 		Body:         payload,
// 		Acknowledger: &MockAcknowledger{},
// 	}
// 	fmt.Printf("%s", payload)
// 	assert.True(t, listener.ListenScoreHistoryBinding(&msg))
// }

func Test_SyncResult(t *testing.T) {
	Init()

	payload2 := request.CreateStudentTarget{
		SmartbtwID:  550094,
		SchoolID:    1,
		MajorID:     1,
		SchoolName:  "POLTEKIM",
		MajorName:   "KELAUTAN",
		TargetScore: 700,
		TargetType:  string(models.PTN),
	}
	_, err1 := lib.CreateStudentTarget(&payload2)
	assert.Nil(t, err1)

	payload := request.CreateHistoryPtn{
		SmartBtwID:              550094,
		TaskID:                  2,
		ModuleCode:              "M-002",
		ModuleType:              string(models.UkaFree),
		PotensiKognitif:         300,
		PenalaranMatematika:     100,
		LiterasiBahasaIndonesia: 120,
		LiterasiBahasaInggris:   122,
		Total:                   500,
		Repeat:                  1,
		ExamName:                "test exam",
		Grade:                   string(models.Basic),
	}
	_, err := lib.CreateHistoryPtn(&payload)
	assert.Nil(t, err)

	payload1 := []byte(`
	{
		"version": 1,
		"data": {
			"smartbtw_id": 550094,
			"task_id": 2,
			"repeat": 1,
			"program": "tps",
			"start" : "2023-01-10T06:39:33.595+00:00",
			"end" : "2023-01-10T06:40:33.595+00:00"
		}
	}
	`)

	msg := amqp.Delivery{
		RoutingKey:   "profile.syncResult",
		Body:         payload1,
		Acknowledger: &MockAcknowledger{},
	}
	assert.True(t, listener.ListenScoreHistoryBinding(&msg))
}

// func Test_UpdateScoreCPNS(t *testing.T) {
// 	Init()

// 	sPhoto := "http://www.google.com"
// 	sPayload := request.BodyMessageCreateUpdateStudent{
// 		ID:    200001,
// 		Name:  "John Doe 2",
// 		Email: "john_doe2@email.com",
// 		Photo: &sPhoto,
// 	}
// 	_, errS := lib.UpsertStudent(&sPayload)
// 	assert.Nil(t, errS)

// 	payload2 := request.CreateWallet{
// 		SmartbtwID: 200001,
// 		Point:      0,
// 		Type:       string(models.BONUS),
// 	}

// 	_, err2 := lib.CreateWallet(&payload2)
// 	if err2 != nil {
// 		assert.NotNil(t, err2)
// 	} else {
// 		assert.Nil(t, err2)
// 	}

// 	py := request.CreateStudentTargetCpns{
// 		SmartbtwID:        200001,
// 		InstanceID:        1,
// 		PositionID:        3,
// 		InstanceName:      "Kementrian Agama",
// 		PositionName:      "testing",
// 		FormationLocation: "PROVINCE",
// 		FormationType:     "CENTRAL",
// 		FormationCode:     "BALADU",
// 		CompetitionID:     1,
// 		TargetScore:       430,
// 	}
// 	_, err := lib.CreateStudentTargetCpns(&py)
// 	assert.Nil(t, err)

// 	payload := []byte(`
// 	{
// 		"version": 2,
// 		"data": {
// 			"smartbtw_id": 200001,
// 			"task_id": 1,
// 			"module_code": "MD-102",
// 			"module_type": "uka_premium",
// 			"tiu": 155,
// 			"tiu_pass_status": true,
// 			"tiu_passing_grade": 80,
// 			"tkp": 194,
// 			"tkp_pass_status": true,
// 			"tkp_passing_grade": 156,
// 			"total": 523,
// 			"twk": 174,
// 			"twk_pass_status": true,
// 			"is_all_passed": false,
// 			"twk_passing_grade": 65,
// 			"repeat": 1,
// 			"exam_name": "GANTENG 2",
// 			"start" : "2023-01-10T06:39:33.595+00:00",
// 			"end" : "2023-01-10T06:40:33.595+00:00",
// 			"is_live": true
// 		}
// 	}
// 	`)

// 	msg := amqp.Delivery{
// 		RoutingKey:   "history-cpns.created",
// 		Body:         payload,
// 		Acknowledger: &MockAcknowledger{},
// 	}
// 	fmt.Printf("%s", payload)
// 	assert.True(t, listener.ListenScoreHistoryBinding(&msg))
// }
