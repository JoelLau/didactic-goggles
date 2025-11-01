package config

import (
	"os"
	"strings"
)

// puts all environment variables into a map
func Env() map[string]string {
	env := make(map[string]string, 0)
	for _, e := range os.Environ() {
		// expects a key-value pairs separated by '=':
		// - tokens[0]: key   (cannot contain '=')
		// - tokens[1]: value (can contain '=')
		tokens := strings.SplitN(e, "=", 2)
		if len(tokens) >= 2 {
			env[tokens[0]] = tokens[1]
			continue
		}
	}

	return env
}
