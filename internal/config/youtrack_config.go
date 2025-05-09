package config

type Youtrack struct {
	BaseURL  string
	ApiToken string
}

func loadYoutrackConfig() *Youtrack {
	return &Youtrack{
		BaseURL:  getEnvWithDefault("YOUTRACK_BASE_URL", "http://localhost:8080/api"),
		ApiToken: getRequiredEnv("YOUTRACK_API_TOKEN"),
	}
}
