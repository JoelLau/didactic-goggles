package config_test

import (
	"didactic-goggles/internal/config"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv(t *testing.T) {
	t.Parallel()

	// [A]rrange
	envKey := fmt.Sprintf("testenv_%s", t.Name()) // TODO: use a uuid prefix
	envVal := "asdf=qwer"

	oldVal := os.Getenv(envKey)
	defer func() { // cleanup
		if oldVal != "" {
			os.Setenv(envKey, oldVal)
			return
		}

		if err := os.Unsetenv(envKey); err != nil {
			t.Fatalf("failed to set env var '%s': %+v", envKey, err)
			return
		}
	}()

	err := os.Setenv(envKey, envVal) // not ideal - avoid modifying env vars where possible
	require.NoErrorf(t, err, "error setting env with key '%s' with value '%+s': %+v", envKey, envVal, err)

	// [A]ct
	env := config.Env()

	// [A]ssert
	got, ok := env[envKey]
	assert.Truef(t, ok, "envKey '%s' must exist in env vars", envKey)
	require.Equal(t, envVal, got)

}
