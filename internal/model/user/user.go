package user

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}
