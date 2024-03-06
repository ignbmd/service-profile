package request

type GetStudentSchoolRanking struct {
	Limit         int    `json:"limit" query:"limit"`
	Page          int    `json:"page" query:"page"`
	SchoolID      string `json:"school_id" query:"school_id"`
	SearchKeyword string `json:"search" query:"search"`
}
