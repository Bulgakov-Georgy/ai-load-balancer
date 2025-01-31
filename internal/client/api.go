package client

import (
	"encoding/json"
	"log"
	"maps"
	"os"

	"ai_load_balancer/internal/configuration"
)

type Api struct {
	Name    string
	Model   string
	Execute func(string, string) (string, error)
	MaxRpm  int
}

var ModelToApis = make(map[string][]Api)

func init() {
	file, err := os.Open(os.Getenv("MODEL_TO_CLIENTS"))
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ModelToApis)
	if err != nil {
		log.Fatal(err)
	}
	config := configuration.Get()
	for model := range maps.Keys(ModelToApis) {
		apis := ModelToApis[model]
		for i := 0; i < len(apis); i++ {
			api := &apis[i]
			switch api.Name {
			case "openai":
				api.Execute = fetchOpenAIResponse
				api.MaxRpm = config.OpenAIRpm
			case "openrouter":
				api.Execute = fetchOpenRouterResponse
				api.MaxRpm = config.OpenRouterRpm
			case "huggingface":
				api.Execute = fetchHuggingFaceResponse
				api.MaxRpm = config.HuggingFaceRpm
			}
		}
	}
}
