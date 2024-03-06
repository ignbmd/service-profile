package mockstruct

type Location struct {
	ID       uint   `json:"_id"`
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
	Type     string `json:"type"`
}
type LocationRequest struct {
	Data []Location `json:"data"`
}
