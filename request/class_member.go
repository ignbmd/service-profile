package request

import (
	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateClassMember struct {
	SmartbtwID  uint               `json:"smartbtw_id"`
	ClassroomID primitive.ObjectID `json:"classroom_id"`
}

type UpdateClassMember struct {
	SmartbtwID        uint               `json:"smartbtw_id"`
	ClassroomIDBefore primitive.ObjectID `json:"classroom_id_before"`
	ClassroomIDAfter  primitive.ObjectID `json:"classroom_id_after"`
}

type ArrayStudentSwitchClass struct {
	SmartbtwID   uint `json:"smartbtw_id"`
	BtwedutechID uint `json:"btwedutech_id"`
}

type SwitchClassMember struct {
	ClassroomID  primitive.ObjectID        `json:"classroom_id"`
	ClassMembers []ArrayStudentSwitchClass `json:"class_members"`
}

type MessageSwitchClassMemberBody struct {
	Version int               `json:"version" valid:"type(int), required"`
	Data    SwitchClassMember `json:"data" valid:"required"`
}
type MessageCreateClassMemberBody struct {
	Version int               `json:"version" valid:"type(int), required"`
	Data    CreateClassMember `json:"data" valid:"required"`
}

type MessageUpdateClassMemberBody struct {
	Version int               `json:"version" valid:"type(int), required"`
	Data    UpdateClassMember `json:"data" valid:"required"`
}

func UnmarshalMessageCreateClassMemberBody(data []byte) (MessageCreateClassMemberBody, error) {
	var decoded MessageCreateClassMemberBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalMessageUpdateClassMemberBody(data []byte) (MessageUpdateClassMemberBody, error) {
	var decoded MessageUpdateClassMemberBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}

func UnmarshalMessageSwitchClassMemberBody(data []byte) (MessageSwitchClassMemberBody, error) {
	var decoded MessageSwitchClassMemberBody
	err := sonic.Unmarshal(data, &decoded)
	return decoded, err
}
