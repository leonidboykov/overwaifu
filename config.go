package overwaifu

// Credentials represents credentials for board login
type Credentials struct {
	Login    string `json:"login" env:"OVERWAIFU_LOGIN,required"`
	Password string `json:"password" env:"OVERWAIFU_PASSWORD,required"`
}
