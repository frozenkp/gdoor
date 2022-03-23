//go:build client

package config

import (
	_ "embed"
)

//go:embed public.key
var Key []byte
