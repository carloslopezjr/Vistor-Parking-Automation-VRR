package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadDotEnv loads key=value pairs from a .env-style file into the process
// environment. Lines starting with '#' and blank lines are ignored. Existing
// variables are not overwritten.
func LoadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		// If the file does not exist, treat it as non-fatal.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Do not override if already set in the environment.
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
	return scanner.Err()
}
