package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GeminiAPIKey string
	Port         string
}

// Load reads the repository's .env while honoring process environment overrides.
func Load() (Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}
	if err := loadRepoEnv(dir); err != nil {
		return Config{}, err
	}

	port, configured := os.LookupEnv("PORT")
	if !configured {
		port = "8080"
	}
	number, err := strconv.ParseUint(port, 10, 16)
	if err != nil || number == 0 {
		return Config{}, errors.New("PORT must be an integer from 1 to 65535")
	}

	return Config{
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		Port:         strconv.FormatUint(number, 10),
	}, nil
}

func loadRepoEnv(dir string) error {
	for {
		// .git can be a directory or a file in a linked worktree.
		_, err := os.Stat(filepath.Join(dir, ".git"))
		if err == nil {
			envPath := filepath.Join(dir, ".env")
			if err := godotenv.Load(envPath); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("could not load repository environment file %s", envPath)
			}
			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
}
