package conf

import (
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

	return &ViperConfig{
		Viper: v,
	}
}
