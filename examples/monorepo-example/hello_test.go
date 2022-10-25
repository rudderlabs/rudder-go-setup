package monorepo_test

import (
	"testing"

	monorepo "github.com/rudderlabs/rudder-go-setup/examples/monorepo-example"
	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	require.Equal(t, "Hello world", monorepo.Hello("world"))
	require.Equal(t, "Hello human", monorepo.Hello("human"))
}
