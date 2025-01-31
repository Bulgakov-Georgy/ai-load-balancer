package configuration

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

// Config holds various service settings from environment variables
type Config struct {
	OpenAIUrl     string
	OpenAIKey     string
	OpenAITimeout time.Duration
	OpenAIRpm     int

	HuggingFaceUrl     string
	HuggingFaceKey     string
	HuggingFaceTimeout time.Duration
	HuggingFaceRpm     int

	OpenRouterUrl     string
	OpenRouterKey     string
	OpenRouterTimeout time.Duration
	OpenRouterRpm     int

	RedisUrl      string
	RedisPassword string
	RedisDB       int

	CacheResponseExpirationTime time.Duration

	UserRpm int
}

var config *Config
var once sync.Once

// Get returns service config
func Get() *Config {
	once.Do(func() {
		config = &Config{
			OpenAIUrl:     os.Getenv("OPENAI_URL"),
			OpenAIKey:     os.Getenv("OPENAI_KEY"),
			OpenAITimeout: getValue(func() (time.Duration, error) { return time.ParseDuration(os.Getenv("OPENAI_TIMEOUT")) }),
			OpenAIRpm:     getValue(func() (int, error) { return strconv.Atoi(os.Getenv("OPENAI_RPM")) }),

			HuggingFaceUrl:     os.Getenv("HUGGINGFACE_URL"),
			HuggingFaceKey:     os.Getenv("HUGGINGFACE_KEY"),
			HuggingFaceTimeout: getValue(func() (time.Duration, error) { return time.ParseDuration(os.Getenv("HUGGINGFACE_TIMEOUT")) }),
			HuggingFaceRpm:     getValue(func() (int, error) { return strconv.Atoi(os.Getenv("HUGGINGFACE_RPM")) }),

			OpenRouterUrl:     os.Getenv("OPENROUTER_URL"),
			OpenRouterKey:     os.Getenv("OPENROUTER_KEY"),
			OpenRouterTimeout: getValue(func() (time.Duration, error) { return time.ParseDuration(os.Getenv("OPENROUTER_TIMEOUT")) }),
			OpenRouterRpm:     getValue(func() (int, error) { return strconv.Atoi(os.Getenv("OPENROUTER_RPM")) }),

			RedisUrl:      os.Getenv("REDIS_URL"),
			RedisPassword: os.Getenv("REDIS_PASSWORD"),
			RedisDB:       0,

			CacheResponseExpirationTime: getValue(func() (time.Duration, error) { return time.ParseDuration(os.Getenv("CACHE_RESPONSE_EXPIRATION_TIME")) }),

			UserRpm: getValue(func() (int, error) { return strconv.Atoi(os.Getenv("USER_RPM")) }),
		}
	})
	return config
}

func getValue[V any](lambda func() (V, error)) V {
	value, err := lambda()
	if err != nil {
		log.Fatal(err)
	}
	return value
}
