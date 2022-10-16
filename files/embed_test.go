package files_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rudderlabs/rudder-go-setup/files"
	"github.com/stretchr/testify/require"
)

func Test_EmbeddedFiles(t *testing.T) {
	dirs, err := files.FS.ReadDir(".")
	require.NoError(t, err)

	var embedFiles []string
	for _, dir := range dirs {
		if dir.IsDir() {
			files, err := files.FS.ReadDir(dir.Name())
			require.NoError(t, err)

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				embedFiles = append(embedFiles, dir.Name()+"/"+file.Name())
			}
		}
	}

	var osFiles []string
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		switch path {
		case "embed.go", "embed_test.go":
			return nil
		}

		if !info.IsDir() {
			osFiles = append(osFiles, path)
		}

		return err
	})
	require.NoError(t, err)

	require.Equal(t, embedFiles, osFiles)
}
