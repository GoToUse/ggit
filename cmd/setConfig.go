package cmd

import "github.com/spf13/viper"

type ConfigYaml struct {
	viper *viper.Viper
}

type (
	GitS struct {
		FilePath  string
		Website   string
		UrlSuffix string
	}
	MirrorUrlS []string
)

func NewConfig() (*ConfigYaml, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("configs/")

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &ConfigYaml{vp}, nil
}

func (c *ConfigYaml) ReadSection(key string, v interface{}) error {
	err := c.viper.UnmarshalKey(key, v)
	if err != nil {
		return err
	}

	return nil
}
