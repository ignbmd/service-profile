package lib

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/olivere/elastic/v7"
	"smartbtw.com/services/profile/db"
	"smartbtw.com/services/profile/mockstruct"
	"smartbtw.com/services/profile/request"
)

func GetStagesCompetitionList(program string, schoolId string, filterType string) ([]request.StagesCompetitionList, error) {

	if filterType == "" {
		filterType = "challenge"
	}

	// Get Stages Challenge Only Data
	stagesData, err := GetStagesChallengeList(program)
	if err != nil {
		return []request.StagesCompetitionList{}, err
	}

	if len(stagesData) < 1 {
		return []request.StagesCompetitionList{}, nil
	}
	stagesList := []request.StagesCompetitionList{}
	ctx := context.Background()
	fetchSize := int64(10000)
	for _, k := range stagesData {
		fmt.Println("stages ", k.Level, " start -> ", time.Now())
		schoolData := map[string][]mockstruct.CreateHistorySimpleData{
			schoolId: {},
		}

		dom := request.StagesCompetitionList{
			Level:   int(k.Level),
			Program: strings.ToLower(program),
			Type:    filterType,
		}

		query := elastic.NewBoolQuery().Must(
			elastic.NewMatchQuery("package_id", k.PackageID),
		)
		totalData, err := db.ElasticClient.Count().
			Index(fmt.Sprintf("student_history_%s", strings.ToLower(program))).
			Query(query).
			Do(ctx)

		if err != nil {
			continue
		}

		if totalData < 1 {
			fmt.Println("no data is available")
			continue
		}

		// chunkData := totalData / fetchSize
		// lastData := totalData - (chunkData * fetchSize)
		fetchProg := 0
		// fmt.Printf("Using chunk %d with last data %d and fetch Size %d\n", chunkData, lastData, fetchSize)

		searchSource := elastic.NewSearchSource()
		searchSource.Query(query)
		sort := elastic.NewFieldSort("_seq_no").Desc()
		searchSource.SortBy(sort)

		searchSource.Size(int(fetchSize)) // set the number of search results to retrieve to 10000

		// create a search service and pass the search source
		searchService := db.ElasticClient.Scroll(fmt.Sprintf("student_history_%s", strings.ToLower(program))).SearchSource(searchSource)

		// scroll through the search results
		for {
			searchResult, err := searchService.Do(context.Background())
			if err != nil {
				if strings.Contains(err.Error(), "EOF") {
					break
				}
				// handle error
				return []request.StagesCompetitionList{}, err
			}
			if searchResult.Hits == nil {
				// no more search results
				break
			}

			// process the search results
			for _, hit := range searchResult.Hits.Hits {
				fetchProg += 1
				t := mockstruct.CreateHistorySimpleData{}
				sonic.Unmarshal(hit.Source, &t)
				// fmt.Printf("Caching %d of %d\n", fetchProg, totalData)

				if t.SchoolOriginID == "" {
					continue
				}
				_, isEx := schoolData[t.SchoolOriginID]
				if !isEx {
					schoolData[t.SchoolOriginID] = []mockstruct.CreateHistorySimpleData{}
				}
				dom.PackageID = t.PackageID
				dom.TaskID = t.TaskID
				schoolData[t.SchoolOriginID] = append(schoolData[t.SchoolOriginID], t)
			}

			searchService.ScrollId(searchResult.ScrollId)
		}

		dom.AttemptedStudent = int(totalData)

		dom.AttemptedStudentOriginSchool = len(schoolData[schoolId])
		dom.AttemptedSchool = len(schoolData)

		stagesList = append(stagesList, dom)
		fmt.Println("stages ", k.Level, " end -> ", time.Now())
	}

	return stagesList, nil
}
