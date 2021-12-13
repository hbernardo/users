package lib

// User represents the user model, contains JSON tags for responses
type User struct {
	ID           string `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	IPAddress    string `json:"ip_address"`
	CreationDate string `json:"creation_date"`
}
