package request

type CreateParentData struct {
	SmartBtwID   int     `json:"smartbtw_id" bson:"smartbtw_id" valid:"type(int),required"`
	ParentName   *string `json:"parent_name" bson:"parent_name" valid:"type(*string),required"`
	ParentNumber *string `json:"parent_number" bson:"parent_number" valid:"type(*string),required"`
}

func (c *CreateParentData) SetParentName(s string) {
	c.ParentName = &s
}

func (c *CreateParentData) SetParentNumber(s string) {
	c.ParentNumber = &s
}
