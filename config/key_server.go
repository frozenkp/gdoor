//go:build server

package config

import (
	_ "embed"
)

//go:embed private.key
var Key []byte
