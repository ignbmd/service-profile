package lib

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/olivere/elastic/v7"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/mockstruct"
)

func GetStudentExplanation(smartbtwID uint, taskID uint, program string) ([]mockstruct.Explanation, error) {
	ctx := context.Background()

	indexName := "answer_explanations"
	if program == "PTN" {
		indexName = "answer_explanations_utbk"
	}
	if program == "CPNS" {
		indexName = "answer_explanations_cpns_skd"
	}

	var t mockstruct.Explanation
	var gres []mockstruct.Explanation

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("smartbtw_id", smartbtwID),
		elastic.NewMatchQuery("legacy_task_id", taskID),
		elastic.NewMatchQuery("repeat", 1),
	)

	res, err := db.ElasticClient.Search().
		Index(indexName).
		Query(query).
		Sort("order", true).
		Size(1000).Do(ctx)

	if err != nil {
		return []mockstruct.Explanation{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []mockstruct.Explanation{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(mockstruct.Explanation))
	}

	return gres, nil
}

func GetStudentExplanationPerQuestion(questionCode string, taskID uint, program string) ([]mockstruct.Explanation, error) {
	ctx := context.Background()

	indexName := "answer_explanations"
	if program == "PTN" {
		indexName = "answer_explanations_utbk"
	}

	if program == "CPNS" {
		indexName = "answer_explanations_cpns_skd"
	}

	var t mockstruct.Explanation
	var gres []mockstruct.Explanation

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("code.keyword", questionCode),
		elastic.NewMatchQuery("legacy_task_id", taskID),
		elastic.NewMatchQuery("repeat", 1),
	)

	res, err := db.ElasticClient.Search().
		Index(indexName).
		Query(query).
		Sort("order", true).
		Size(1000).Do(ctx)

	if err != nil {
		return []mockstruct.Explanation{}, err
	}

	recCount := res.TotalHits()

	if recCount == 0 {
		return []mockstruct.Explanation{}, nil
	}

	for _, item := range res.Each(reflect.TypeOf(t)) {
		gres = append(gres, item.(mockstruct.Explanation))
	}

	return gres, nil
}

func GetSingleQuestionFromMaster(code string, program string) (*mockstruct.QuestionsElastic, error) {
	question := mockstruct.QuestionsElastic{}
	ctx := context.Background()

	qstIds := strings.Split(code, "-")
	if len(qstIds) < 2 {
		return nil, errors.New("^invalid question code")
	}

	index := "master_questions"
	if program == "PTN" {
		index = "master_questions_utbk"
	}

	query := elastic.NewBoolQuery().Must(
		elastic.NewMatchQuery("id.keyword", qstIds[1]),
		// elastic.NewMatchQuery("group.keyword", "GROUP-1"),
	)
	sort := elastic.NewFieldSort("_seq_no").Desc()

	res, err := db.ElasticClient.Search().
		Index(index).
		SortBy(sort).
		Query(query).
		Size(1).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var t mockstruct.QuestionsElastic
	for _, item := range res.Each(reflect.TypeOf(t)) {
		question = item.(mockstruct.QuestionsElastic)

	}
	return &question, nil
}
