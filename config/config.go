package config

import (
	"os"
	"path"

	"github.com/pelletier/go-toml"
)

// DefaultPath returns the default path for the config file.
func DefaultPath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homedir, ".config", "produgit"), nil
}

// DefaultConfigPath returns the default path for the config file.
func DefaultConfigPath() (string, error) {
	configPath, err := DefaultPath()
	if err != nil {
		return "", err
	}
	return path.Join(configPath, "config.toml"), nil
}

// DefaultOutputPath returns the default path for the output file.
func DefaultOutputPath() (string, error) {
	configPath, err := DefaultPath()
	if err != nil {
		return "", err
	}
	return path.Join(configPath, "report.pb"), nil
}

// Report is the configuration for the report command.
type Report struct {
	Exclude []string `toml:"exclude"`
	Output  string   `toml:"output"`
}

// Plot is the configuration for the plot command.
type Plot struct {
	Output  string   `toml:"output"`
	Authors []string `toml:"authors"`
}

// Config is the configuration for the produgit command.
type Config struct {
	Report *Report `toml:"report"`
	Plot   *Plot   `toml:"plot"`
}

// New creates a new Config with default values.
func New() (*Config, error) {
	defaultOutputPath, err := DefaultOutputPath()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Report: &Report{
			Exclude: []string{
				"**node_modules/*",
				"**vendor/*",
				"**go.sum",
				"**go.mod",
				"**yarn.lock",
				"**package-lock.json",
				"**requirements.txt",
				"**venv/*",
				"**pnpm-lock.yaml",
				"**dist/*",
				"**build/*",
				"**.git/*",
				"**.idea/*",
				"**.vscode/*",
				"**.pytest_cache/*",
				"**.next/*",
				"**.cache/*",
				"**__pycache__/*",
				"**coverage.xml",
				"**coverage.html",
				"**coverage.lcov",
				"**LICENSE.md",
				"*.csv",
				"*.pdf",
				"*.doc",
				"*.docx",
				"*.json",
				"*.png",
				"*.jpg",
				"*.jpeg",
				"*.gif",
				"*.svg",
				"*.ico",
				"*.woff",
				"*.woff2",
				"*.ttf",
				"*.eot",
				"*.txt",
				"**.DS_Store",
				"**Thumbs.db",
				"*.log",
				"*.bak",
				"*.swp",
				"*.swo",
				"*.tmp",
				"*.temp",
				"*.o",
				"*.out",
				"*.jar",
				"*.war",
				"*.ear",
				"*.sqlite3",
			},
			Output: defaultOutputPath,
		},
		Plot: &Plot{
			Output:  "<chart>_<authors>_<date>.html",
			Authors: []string{},
		},
	}

	return cfg, nil
}

// Load loads a config file from the given path.
func Load() (*Config, error) {
	defaultPath, err := DefaultConfigPath()
	if err != nil {
		return nil, err
	}

	if !Exists() {
		if err := Reset(); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(defaultPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	if err := toml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// Reset creates a new default config file.
func Reset() error {
	if err := Remove(); err != nil {
		return err
	}

	defaultConfigPath, err := DefaultConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(defaultConfigPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(defaultConfigPath)
	if err != nil {
		return nil
	}
	defer f.Close()

	cfg, err := New()
	if err != nil {
		return err
	}

	return toml.NewEncoder(f).Encode(cfg)
}

// Exists returns true if the default config file exists.
func Exists() bool {
	defaultPath, err := DefaultConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(defaultPath)
	return !os.IsNotExist(err)
}

// Remove removes the default config file.
func Remove() error {
	if !Exists() {
		return nil
	}

	defaultPath, err := DefaultConfigPath()
	if err != nil {
		return err
	}

	return os.Remove(defaultPath)
}
