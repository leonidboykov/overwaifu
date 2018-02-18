package overwaifu

// SankakuCredentials represents credentials for board login
type SankakuCredentials struct {
	User     string `json:"user" env:"OVERWAIFU_SANKAKU_USER,required"`
	Password string `json:"password" env:"OVERWAIFU_SANKAKU_PASS,required"`
}

// MongoDBConfig ...
type MongoDBConfig struct {
	User           string   `json:"user" env:"OVERWAIFU_MGO_USER,required"`
	Password       string   `json:"password" env:"OVERWAIFU_MGO_PASS,required"`
	URI            []string `json:"uri" env:"OVERWAIFU_MGO_URI,required" envSeparator:","`
	Source         string   `json:"source" env:"OVERWAIFU_MGO_SOURCE,required"`
	ReplicaSetName string   `json:"replica" env:"OVERWAIFU_MGO_REPLICA,required"`
}
