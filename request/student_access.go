package request

type CreateStudentAccess struct {
	SmartBtwID       int    `json:"smartbtw_id" valid:"required"`
	DisallowedAccess string `json:"disallowed_access" valid:"required"`
	AppType          string `json:"app_type" valid:"required"`
}

type CreateStudentAccessBulk struct {
	SmartBtwID       int      `json:"smartbtw_id" valid:"required"`
	DisallowedAccess []string `json:"disallowed_access" valid:"required"`
	AppType          string   `json:"app_type" valid:"required"`
}

type UpdateStudentAccess struct {
	DisallowedAccess string `json:"disallowed_access"`
	AppType          string `json:"app_type"`
}

type DeleteStudentAccess struct {
	DisallowedAccess     string   `json:"disallowed_access"`
	DisallowedAccessBulk []string `json:"disallowed_access_bulk"`
	AppType              string   `json:"app_type"`
}

type GetStudentAccessElastic struct {
	Code    string `json:"code" query:"code"`
	AppType string `json:"app_type" query:"app_type"`
}
