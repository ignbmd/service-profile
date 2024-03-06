package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	b, err := sonic.Marshal(obj)
	if err != nil {
		return map[string]interface{}{}, err
	}
	var m map[string]interface{}
	err = sonic.Unmarshal(b, &m)
	if err != nil {
		return map[string]interface{}{}, err
	}

	fin := map[string]interface{}{}
	for n, k := range m {
		if k != nil {
			fin[n] = k
		}
	}
	return fin, nil
}

func AppendCategory(a []string, b []string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter := range check {
		res = append(res, letter)
	}

	return res
}

func RemoveDuplicates(arr []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range arr {
		cleanedString := strings.ReplaceAll(v, ",", "")
		cleanedString = strings.ReplaceAll(cleanedString, ".", "")

		if encountered[cleanedString] == false {
			encountered[cleanedString] = true
			result = append(result, cleanedString)
		}
	}

	return result
}

func ToLowerAndUnderscore(source string) string {
	return strings.ReplaceAll(strings.ToLower(source), " ", "_")
}

func RemoveUnderscore(source string) string {
	return strings.ReplaceAll(strings.ToLower(source), "_", " ")
}

type SubCategories struct {
	DateStart     string `json:"date_start"`
	ExamType      string `json:"exam_type"`
	Materi        string `json:"materi"`
	Score         int    `json:"score"`
	ScoreMax      int    `json:"score_max"`
	TaskID        uint   `json:"task_id"`
	Repeat        int    `json:"repeat"`
	FormattedDate string `json:"formatted_date"`
	FormattedTime string `json:"formatted_time"`
}

type DataStruct struct {
	Category      string        `json:"category"`
	CategoryName  string        `json:"category_name"`
	Program       string        `json:"program"`
	TaskID        uint          `json:"task_id"`
	Repeat        int           `json:"repeat"`
	SubCategories SubCategories `json:"sub_categories"`
}

func FormatingCategoryPTN(data []string) (map[string]SubCategories, error) {
	finalResults := []SubCategories{}
	taskIDMap := make(map[uint]SubCategories)

	for _, member := range data {
		var data SubCategories
		if err := sonic.Unmarshal([]byte(member), &data); err != nil {
			continue
		}

		if existingData, ok := taskIDMap[data.TaskID]; ok {
			if data.Repeat > existingData.Repeat {
				taskIDMap[data.TaskID] = data
			}
		} else {
			taskIDMap[data.TaskID] = data
		}
	}

	for _, value := range taskIDMap {
		finalResults = append(finalResults, value)
	}

	resultsFormatted := map[string]SubCategories{}
	for _, pre := range finalResults {
		parsedTime, err := ParseTime(pre.DateStart)
		if err != nil {
			continue
		}
		formattedDate := parsedTime.Format("02/01/2006")
		formattedTime := parsedTime.Format("15:04:05")
		pre.FormattedDate = formattedDate
		pre.FormattedTime = formattedTime

		// if strings.Contains(strings.ToLower(pre.Materi), "pengetahuan kuantitatif") {
		// 	pre.Materi = "utbk pk"
		// 	resultsFormatted[ToLowerAndUnderscore(pre.Materi)] = pre
		// }
		resultsFormatted[ToLowerAndUnderscore(pre.Materi)] = pre
	}

	return resultsFormatted, nil
}

func FormatingCategory(data []string) (map[string]DataStruct, error) {

	finalResults := []DataStruct{}
	taskIDMap := make(map[uint]DataStruct)

	for _, member := range data {
		var data DataStruct
		if err := sonic.Unmarshal([]byte(member), &data); err != nil {
			continue
		}
		// fmt.Println(data)
		if existingData, ok := taskIDMap[data.TaskID]; ok {
			if data.Repeat > existingData.Repeat {
				taskIDMap[data.TaskID] = data
			}
		} else {
			taskIDMap[data.TaskID] = data
		}
	}

	for _, value := range taskIDMap {
		finalResults = append(finalResults, value)
	}

	resultsFormatted := map[string]DataStruct{}
	for _, pre := range finalResults {
		parsedTime, err := ParseTime(pre.SubCategories.DateStart)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return nil, err
		}
		formattedDate := parsedTime.Format("02/01/2006")
		formattedTime := parsedTime.Format("15:04:05")
		pre.SubCategories.FormattedDate = formattedDate
		pre.SubCategories.FormattedTime = formattedTime

		if strings.Contains(RemoveUnderscore(pre.SubCategories.Materi), "nasionalisme") || strings.Contains(RemoveUnderscore(pre.SubCategories.Materi), "skd nasionalisme") {
			pre.Category = "TWK"
			pre.SubCategories.Materi = "nasionalisme"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "integritas") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd integritas") {
			pre.Category = "TWK"
			pre.SubCategories.Materi = "integritas"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "bela negara") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd bela negara") {
			pre.Category = "TWK"
			pre.SubCategories.Materi = "bela negara"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "bahasa indonesia") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd bahasa indonesia") {
			pre.Category = "TWK"
			pre.SubCategories.Materi = "bahasa indonesia"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "pilar negara") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd pilar negara") {
			pre.Category = "TWK"
			pre.SubCategories.Materi = "pilar negara"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "verbal silogisme") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd verbal silogisme") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "verbal silogisme"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "figural ketidaksamaan gamber") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd figural ketidaksamaan gamber") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "figural ketidaksamaan gamber"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "figural serial gambar") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd figural serial gambar") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "figural serial gambar"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "figural analogi") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd figural analogi") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "figural analogi"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "numerik berhitung") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd numerik berhitung") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "numerik berhitung"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "numerik perbandingan") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd numerik perbandingan") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "numerik perbandingan"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "verbal analitis") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd verbal analitis") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "verbal analitis"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "verbal analogi") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd verbal analogi") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "verbal analogi"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "numerik soal cerita") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd numerik soal cerita") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "numerik soal cerita"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "numerik deret") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd numerik deret") {
			pre.Category = "TIU"
			pre.SubCategories.Materi = "numerik deret"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "jejaring kerja") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd jejaring kerja") {
			pre.Category = "TKP"
			pre.SubCategories.Materi = "jejaring kerja"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "profesionalisme") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd profesionalisme") {
			pre.Category = "TKP"
			pre.SubCategories.Materi = "profesionalisme"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "sosial budaya") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd sosial budaya") {
			pre.Category = "TKP"
			pre.SubCategories.Materi = "sosial budaya"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(RemoveUnderscore(pre.SubCategories.Materi), "teknologi informasi dan komunikasi") || strings.Contains(RemoveUnderscore(pre.SubCategories.Materi), "skd tik") {
			pre.Category = "TKP"
			pre.SubCategories.Materi = "teknologi informasi dan komunikasi"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "anti radikalisme") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd anti radikalisme") {
			pre.Category = "TKP"
			pre.SubCategories.Materi = "anti radikalisme"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}
		if strings.Contains(strings.ToLower(pre.SubCategories.Materi), "pelayanan publik") || strings.Contains(strings.ToLower(RemoveUnderscore(pre.SubCategories.Materi)), "skd pelayanan publik") {
			pre.Category = "TKP"
			pre.SubCategories.Materi = "pelayanan publik"
			resultsFormatted[RemoveUnderscore(pre.SubCategories.Materi)] = pre
		}

	}

	return resultsFormatted, nil
}

func ParseTime(timeString string) (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04:05.999999-07:00",
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.Parse(format, timeString)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func ConvertToTitleCase(input string) string {
	words := strings.Split(input, "_")

	for i, word := range words {
		words[i] = strings.Title(word)
	}

	result := strings.Join(words, " ")

	return result
}

func ArrayContainsCS(s []string, searchterm string) bool {
	for _, k := range s {
		if strings.Contains(strings.ToLower(k), strings.ToLower(searchterm)) {
			return true
		}
	}
	return false
}
