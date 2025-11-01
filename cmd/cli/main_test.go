package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	t.Parallel()

	app := App{
		Environment:  map[string]string{},
		InputStream:  os.Stdin,
		ErrorStream:  os.Stderr,
		OutputStream: os.Stdout,
	}

	err := app.Run(t.Context(), []string{"../../tests/testdata/dbs_statement.csv"})
	require.NoError(t, err)
}

func TestDSN(t *testing.T) {
	given := map[string]string{
		"DB_HOST": "192.168.0.1",
		"DB_PORT": "5432",
		"DB_USER": "ddgadmin",
		"DB_PASS": "some_secure_p@ssw0rd??",
		"DB_NAME": "ddg",
	}
	want := "postgres://ddgadmin:some_secure_p@ssw0rd??@192.168.0.1:5432/ddg"

	app := &App{Environment: given}
	got := app.Dsn()
	require.Equal(t, got, want)
}
