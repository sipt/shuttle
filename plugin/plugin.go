package plugin

import (
	"os"
	"path"
	"plugin"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sirupsen/logrus"
)

const (
	PluginsDir = "plugins"
	ExtName    = ".plugin"
)

func ApplyConfig(config *model.Config) error {
	if dir, err := os.Stat(PluginsDir); err != nil && !os.IsExist(err) {
		logrus.Infof("not found plugins")
		return nil
	} else if !dir.IsDir() {
		logrus.Errorf("plugins is not a dir")
	}
	dir, err := os.Open(PluginsDir)
	if err != nil {
		return errors.Wrap(err, "open plugins dir failed")
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return errors.Wrap(err, "get plugins failed")
	}
	for _, file := range files {
		if !file.IsDir() && path.Ext(file.Name()) == ExtName {
			_, name := path.Split(file.Name())
			name = name[:len(name)-len(ExtName)]
			pluginPath := path.Join(PluginsDir, file.Name())
			plug, err := plugin.Open(pluginPath)
			if err != nil {
				logrus.WithError(err).Errorf("import plugin[%s] failed", pluginPath)
				continue
			}
			f, err := plug.Lookup("ApplyConfig")
			if err != nil {
				logrus.WithError(err).Errorf("Lookup [func ApplyConfig(*model.Config) error] in plugin [%s] failed", file.Name())
				continue
			}
			applyConfig, ok := f.(func(map[string]string) error)
			if ok {
				err = applyConfig(config.Plugins[name])
				if err != nil {
					logrus.WithError(err).Errorf("plugin[%s].ApplyConfig failed", name)
					continue
				}
			}
		}
	}
	return nil
}
