package files

import "embed"

//go:embed **/*.yaml
//go:embed **/*.yaml.tmpl
//go:embed **/Makefile
var FS embed.FS
