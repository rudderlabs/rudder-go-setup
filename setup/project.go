package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/rudderlabs/rudder-go-setup/files"
	"golang.org/x/mod/modfile"
)

type Project struct {
	Name           string
	GoVersion      string
	ProjectPath    string
	RepositoryRoot string
	ProjectRelPath string // relative path to repo root
}

func (p *Project) Detect() error {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("failed to detect repository root: %w", err)
	}
	p.RepositoryRoot = path.Base(strings.TrimSpace(string(out)))

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	p.ProjectPath = wd
	modData, err := os.ReadFile(path.Join(p.ProjectPath, "go.mod"))
	if err != nil {
		return fmt.Errorf("open go.mod in working directory: %w", err)
	}

	mfile, err := modfile.Parse(path.Join(wd, "go.mod"), modData, nil)
	if err != nil {
		return fmt.Errorf("failed to parse go.mod in working directory: %w", err)
	}

	p.ProjectRelPath = path.Clean(strings.TrimPrefix(wd, p.RepositoryRoot))

	p.GoVersion = mfile.Go.Version
	p.Name = path.Base(mfile.Module.Mod.Path)

	return nil
}

func (p *Project) Init() error {
	pt, err := template.ParseFS(files.FS, "project/**")
	if err != nil {
		return err
	}

	projectFiles := map[string]string{
		"Makefile":      "Makefile",
		"golangci.yaml": ".golangci.yaml",
	}

	for src, dest := range projectFiles {
		fmt.Printf("Adding %s \n", dest)

		f, err := os.OpenFile(path.Join(p.ProjectPath, dest), os.O_RDWR|os.O_CREATE, 0o666)
		if err != nil {
			return err
		}

		if err := pt.ExecuteTemplate(f, src, p); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	repoFiles := map[string]string{
		"test.yaml.tmpl": ".github/workflows/test.yaml",
	}

	for src, dest := range repoFiles {
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(dest), 0700)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Adding %s \n", dest)

		f, err := os.OpenFile(path.Join(p.ProjectPath, dest), os.O_RDWR|os.O_CREATE, 0o666)
		if err != nil {
			return err
		}

		if err := pt.ExecuteTemplate(f, src, p); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}
