package config

type GlobalConfig struct {
	AppName string
}

var LoadedGlobalConfig *GlobalConfig

func loadGlobalConfig() {
	LoadedGlobalConfig = &GlobalConfig{
		AppName: getEnvWithDefault("APP_NAME", "Catfeinated Time Tracker"),
	}
}
