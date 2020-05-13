package Config

var myConfig *INIConfig

func InitConfig(confName string) bool {
	myConfig = new(INIConfig)
	myConfig.InitConfig(confName)
	return true
}

func ReadKey(path string, key string) string {
	return myConfig.Read(path, key)
}
