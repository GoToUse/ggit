package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

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
	homePath, _ := os.UserHomeDir()
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	currentConfPath, _ := filepath.Abs("configs/")
	ggitDir := filepath.Join(homePath, ".ggit")

	_, err := os.Stat(currentConfPath)
	_, errGgit := os.Stat(ggitDir)
	if os.IsNotExist(err) {
		if os.IsNotExist(errGgit) {
			return nil, errors.New(fmt.Sprintf("%s or %s not exist.", currentConfPath, ggitDir))
		} else {
			vp.AddConfigPath(ggitDir + "/")
			err = vp.ReadInConfig()
			if err != nil {
				return nil, err
			}
			return &ConfigYaml{vp}, nil
		}
	} else {
		vp.AddConfigPath("configs/")
		err = vp.ReadInConfig()
		if err != nil {
			return nil, err
		}
		return &ConfigYaml{vp}, nil
	}
}

func (c *ConfigYaml) ReadSection(key string, v interface{}) error {
	err := c.viper.UnmarshalKey(key, v)
	if err != nil {
		return err
	}

	return nil
}
