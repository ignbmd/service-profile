package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivere/elastic/v7"
	"smartbtw.com/services/profile/config"
)

var ElasticClient *elastic.Client

func NewElastic() {
	if ElasticClient != nil {
		return
	}
	urls := strings.Split(os.Getenv("ELASTICSEARCH_URLS"), ",")
	username := os.Getenv("ELASTICSEARCH_USERNAME")
	password := os.Getenv("ELASTICSEARCH_PASSWORD")
	if len(urls) == 0 {
		urls = append(urls, "http://localhost:9200")
	}
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(urls...),
		elastic.SetBasicAuth(username, password),
	)
	if err != nil {
		log.Fatal(err)
	}
	ElasticClient = client
}

func GetStudentTargetPtkIndexName() (name string) {
	name = "student_profile_ptk"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetStudentTargetPtnIndexName() (name string) {
	name = "student_profile_ptn"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetStudentTargetCpnsIndexName() (name string) {
	name = "student_profile_cpns"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

// func GetStudentPtnProfileIndexName() (name string) {
// 	name = "student_profile_ptn"
// 	if config.IsTest() {
// 		name = fmt.Sprintf("test_%s", name)
// 	}
// 	return
// }

// func GetStudentPtkProfileIndexName() (name string) {
// 	name = "student_profile_ptk"
// 	if config.IsTest() {
// 		name = fmt.Sprintf("test_%s", name)
// 	}
// 	return
// }

func GetStudentHistoryPtnIndexName() (name string) {
	name = "student_history_ptn"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetClassMemberIndexName() (name string) {
	name = "class_members"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetStudentHistoryPtkIndexName() (name string) {
	name = "student_history_ptk"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetStudentHistoryAssessmentIndexName(program string) (name string) {
	name = fmt.Sprintf("student_history_%s", program)
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetStudentHistoryCpnsIndexName() (name string) {
	name = "student_history_cpns"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetStudentProfileIndexName() (name string) {
	name = "student_profiles"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetStudentDisallowedAccessIndexName() (name string) {
	name = "student_disallowed_access"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func GetProfileCPNSExamResult() (name string) {
	name = "skd_cpns_exam_result"
	if config.IsTest() {
		name = fmt.Sprintf("test_%s", name)
	}
	return
}

func setup(indexes []string) {
	ctx := context.Background()
	for _, v := range indexes {
		ElasticClient.CreateIndex(v).Do(ctx)
	}
}

func teardown(indexes []string) {
	ctx := context.Background()
	for _, v := range indexes {
		ElasticClient.DeleteIndex(v).Do(ctx)
	}
}

func WithElasticSetupTeardown(fn func()) (err error) {
	if err = os.Setenv("ENV", config.EnvTest); err != nil {
		return
	}
	indexes := []string{
		GetStudentTargetPtkIndexName(),
		GetStudentTargetPtnIndexName(),
		// GetStudentPtkProfileIndexName(),
		// GetStudentPtnProfileIndexName(),
		GetStudentHistoryPtkIndexName(),
		GetStudentHistoryPtnIndexName(),
		GetStudentHistoryCpnsIndexName(),
		GetStudentTargetCpnsIndexName(),
	}
	teardown(indexes)
	setup(indexes)
	fn()
	// teardown(indexes)
	return
}
