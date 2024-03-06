package mockstruct

type Stages struct {
	Stage         uint   `json:"stage" bson:"stage"`
	Level         uint   `json:"level" bson:"level"`
	RequiredStage uint   `json:"required_stage" bson:"required_stage"`
	ModuleType    string `json:"module_type" bson:"module_type"`
	PackageID     uint   `json:"package_id" bson:"package_id"`
	Type          string `json:"type" bson:"type"`
	ProductCode   string `json:"product_code" bson:"product_code"`
	IsLocked      bool   `json:"is_locked" bson:"is_locked"`
	Session       string `json:"session" bson:"session"`
	IsScheduled   bool   `json:"is_scheduled" bson:"is_scheduled"`
}
