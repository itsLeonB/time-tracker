package config

type GlobalConfig struct {
	AppName string
}

var LoadedGlobalConfig *GlobalConfig

func init() {
	LoadedGlobalConfig = &GlobalConfig{
		AppName: "Time Tracker",
	}
}
