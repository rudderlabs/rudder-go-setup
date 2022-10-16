package monorepo_test

import (
	"testing"

	"github.com/rudderlabs/rudder-go-setup/examples/monorepo"
	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	require.Equal(t, "Hello world", monorepo.Hello("world"))
	require.Equal(t, "Hello human", monorepo.Hello("human"))
}
