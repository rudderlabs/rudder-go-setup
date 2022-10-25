package files

import "embed"

//go:embed **/*.yaml
//go:embed **/*.yaml.tmpl
//go:embed **/Makefile
//go:embed **/Dockerfile
var FS embed.FS
