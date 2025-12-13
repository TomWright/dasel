package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"go.yaml.in/yaml/v4"
)

// Config holds the contents of a config file.
type Config struct {
	DefaultFormat string `yaml:"default_format"`
}

var cfg = Config{
	DefaultFormat: "json",
}
var cfgLoaded = false

// LoadConfig loads the config from the given path.
// If already loaded, returned previously loaded config.
func LoadConfig(path string) (Config, error) {
	if cfgLoaded {
		return cfg, nil
	}

	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return cfg, fmt.Errorf("error getting current user: %v", err)
		}
		path = usr.HomeDir + path[1:]
	}

	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		cfgLoaded = true
		return cfg, nil
	}
	if err != nil {
		return cfg, fmt.Errorf("error opening config file at path %q: %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			// Noop
		}
	}()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil && !errors.Is(err, io.EOF) {
		return cfg, fmt.Errorf("error parsing config file: %w", err)
	}
	cfgLoaded = true
	return cfg, nil
}
