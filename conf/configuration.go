package conf

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Configuration holds configuration struct
type Configuration struct {
	SC struct {
		Username string `required:"true"`
		Password string `required:"true"`
	}
	DB struct {
		Username       string `required:"true"`
		Password       string `required:"true"`
		URI            string `required:"true"`
		Source         string `required:"true"`
		ReplicaSetName string `split_words:"true" required:"true"`
	}
	Netlify struct {
		BuildHook string `split_words:"true" required:"true"`
	}
	// JSONBin struct {
	// 	APIKey string `envconfig:"api_key" required:"true"`
	// }
	MyJSON struct {
		BucketID string `split_words:"true" required:"true"`
	}
}

func loadEnvironment(filename string) error {
	if filename != "" {
		return godotenv.Load(filename)
	}
	err := godotenv.Load()
	// handle error is .env does not exist
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Load loads configuration
func Load(filename string) (*Configuration, error) {
	if err := loadEnvironment(filename); err != nil {
		return nil, err
	}

	config := new(Configuration)
	if err := envconfig.Process("overwaifu", config); err != nil {
		return nil, err
	}

	return config, nil
}
