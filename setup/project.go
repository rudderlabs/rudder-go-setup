package setup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/rudderlabs/rudder-go-setup/files"
	"golang.org/x/mod/modfile"
)

var emptyComment = regexp.MustCompile(`^\s*#\s*$`)

type Project struct {
	Name           string
	GoVersion      string
	ProjectPath    string
	RepositoryRoot string
	ProjectRelPath string // relative path to repo root
	Nested         bool
}

func (p *Project) Detect() error {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return fmt.Errorf("failed to detect repository root: %w", err)
	}
	p.RepositoryRoot = path.Clean(strings.TrimSpace(string(out)))

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

	p.ProjectRelPath = strings.TrimPrefix(strings.TrimPrefix(wd, p.RepositoryRoot), "/")
	if p.ProjectRelPath != "" {
		p.Nested = true
	}

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
		"Dockerfile":    "Dockerfile",
		"golangci.yaml": ".golangci.yaml",
	}

	for src, dest := range projectFiles {
		fmt.Printf("Adding %s \n", dest)

		f, err := os.Create(path.Join(p.ProjectPath, dest))
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
		"test.yaml.tmpl": fmt.Sprintf(".github/workflows/test-%s.yaml", p.Name),
		"build.yaml":     fmt.Sprintf(".github/workflows/build-%s.yaml", p.Name),
	}

	for src, dest := range repoFiles {
		destCanonical := path.Join(p.RepositoryRoot, dest)

		if _, err := os.Stat(destCanonical); os.IsNotExist(err) {
			err = os.MkdirAll(path.Dir(destCanonical), 0o700)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Adding %s \n", dest)

		f, err := os.Create(destCanonical)
		if err != nil {
			return err
		}
		// filter for empty comments, template artifacts
		ff := filter(f, emptyComment.Match)

		if err := pt.ExecuteTemplate(ff, src, p); err != nil {
			return err
		}

		if err := ff.Close(); err != nil {
			return err
		}
	}
	return nil
}

func filter(output io.Writer, fn func(line []byte) bool) io.WriteCloser {
	r, w := io.Pipe()

	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			if fn(s.Bytes()) {
				continue
			}

			_, err := output.Write(s.Bytes())
			if err != nil {
				w.CloseWithError(err)
				return
			}

			_, err = output.Write([]byte{'\n'})
			if err != nil {
				w.CloseWithError(err)
				return
			}
		}
		w.CloseWithError(s.Err())
	}()
	return w
}
