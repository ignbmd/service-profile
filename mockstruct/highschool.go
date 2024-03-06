package mockstruct

type SchoolLocation struct {
	ID         string `json:"_id"`
	Name       string `json:"name"`
	LocationID uint   `json:"location_id"`
	Type       string `json:"type"`
}
type SchoolLocationRequest struct {
	Data SchoolLocation `json:"data"`
}
