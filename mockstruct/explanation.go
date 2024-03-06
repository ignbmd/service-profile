package mockstruct

import "time"

type Explanation struct {
	ID                string   `json:"id"`
	Code              string   `json:"code"`
	Group             string   `json:"group"`
	Module            string   `json:"module"`
	SmartbtwID        int      `json:"smartbtw_id"`
	Repeat            int      `json:"repeat"`
	Order             int      `json:"order"`
	Question          *string  `json:"question"`
	Explanation       string   `json:"explanation"`
	Category          string   `json:"category"`
	SubCategory       *string  `json:"sub_category"`
	OptionKeys        []int    `json:"option_keys"`
	OptionTypes       []string `json:"options_types"`
	OptionID          []int    `json:"option_id"`
	OptionTexts       []string `json:"option_texts"`
	OptionValues      []int    `json:"option_values"`
	Answered          int      `json:"answered_option_id"`
	AnswerKey         int      `json:"answer_key_id"`
	CreatedAt         string   `json:"created_at"`
	LegacyTaskID      int      `json:"legacy_task_id"`
	IsUnderstand      *bool    `json:"is_understand"`
	TimeConsumed      int      `json:"time_consumed"`
	AnsweredTrueItem  []uint   `json:"answered_true_list"`
	AnsweredFalseItem []uint   `json:"answered_false_list"`
	AnswerType        string   `json:"answer_type"`
	AnswerHeaderTrue  *string  `json:"answer_header_true"`
	AnswerHeaderFalse *string  `json:"answer_header_false"`
	Essay             *string  `json:"essay"`
	AnsweredEssay     *string  `json:"answered_essay"`
}

type QuestionsElastic struct {
	ID                string    `json:"id"`
	Question          string    `json:"question"`
	QuestionType      string    `json:"question_type"`
	AnswerType        string    `json:"answer_type"`
	Category          string    `json:"category"`
	CategoryID        uint      `json:"category_id"`
	SubCategory       string    `json:"sub_category"`
	SubCategoryID     uint      `json:"sub_category_id"`
	Explanation       string    `json:"explanation"`
	ExplanationMedia  *string   `json:"explanation_media"`
	Program           string    `json:"program"`
	ParentID          *uint     `json:"parent_id"`
	OptionTypes       []string  `json:"option_types"`
	OptionIDs         []uint    `json:"option_ids"`
	OptionTexts       []string  `json:"option_texts"`
	OptionValues      []float64 `json:"option_values"`
	QuestionKeyword   []string  `json:"question_keyword"`
	AnswerHeaderTrue  *string   `json:"answer_header_true"`
	AnswerHeaderFalse *string   `json:"answer_header_false"`
	AnswerEssay       *string   `json:"answer_essay"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
