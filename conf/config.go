package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

// ViperConfig ...
type ViperConfig struct {
	*viper.Viper
}

// StudyWS ...
var StudyWS *ViperConfig

func init() {
	StudyWS = readConfig(map[string]interface{}{
		"debug_route": false,
		"port":        10812,
	})
}

func readConfig(defaults map[string]interface{}) *ViperConfig {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	v.AddConfigPath("./conf")

	v.AutomaticEnv()

	v.SetConfigName(".env.dev")

	err := v.ReadInConfig()
	if err != nil {
		fmt.Println("conf", "readConfig", "Error", err) // TODO: logger
		return nil
	}

	return &ViperConfig{
		Viper: v,
	}
}
