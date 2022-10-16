package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	main "github.com/rudderlabs/rudder-go-setup"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	require.NoError(t, os.Chdir("./examples/monorepo"))

	err := main.App.Run([]string{"go-setup", "init"})
	require.NoError(t, err)

	makeArgs := []string{"test", "fmt", "lint"}

	for _, arg := range makeArgs {
		t.Run(fmt.Sprintf("make %s", arg), func(t *testing.T) {
			out, err := exec.Command("make", arg).CombinedOutput()
			t.Logf("%s\n", out)
			require.NoError(t, err)
		})
	}
}
