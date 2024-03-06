package request

import "time"

type StudentModuleProgress struct {
	SmartBtwID  int         `json:"smartbtw_id" bson:"smartbtw_id" valid:"type(int),required"`
	TaskID      int         `json:"task_id" bson:"task_id" valid:"type(int),required"`
	ModuleNo    int         `json:"module_no" bson:"module_no" valid:"type(int),required"`
	Repeat      int         `json:"repeat" bson:"repeat" valid:"type(int),required"`
	ModuleTotal int         `json:"module_total" bson:"module_total" valid:"type(int),required"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at" valid:"type(time.Time)"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at" valid:"type(time.Time)"`
	DeletedAt   interface{} `json:"deleted_at" bson:"deleted_at" valid:"optional"`
}

type UpdateStudentModuleProgress struct {
	TaskID      int         `json:"task_id" bson:"task_id" valid:"type(int),required"`
	ModuleNo    int         `json:"module_no" bson:"module_no" valid:"type(int),required"`
	Repeat      int         `json:"repeat" bson:"repeat" valid:"type(int),required"`
	ModuleTotal int         `json:"module_total" bson:"module_total" valid:"type(int),required"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at" valid:"type(time.Time)"`
	UpdatedAt   time.Time   `json:"updated_at" bson:"updated_at" valid:"type(time.Time)"`
	DeletedAt   interface{} `json:"deleted_at" bson:"deleted_at" valid:"optional"`
}
