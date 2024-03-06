package request

type InterviewSessionRequest struct {
	Name        string `json:"name" bson:"name" valid:"required"`
	Description string `json:"description" bson:"description" valid:"required"`
	Number      int    `json:"number" bson:"number" valid:"required"`
}
