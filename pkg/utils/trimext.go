package utils

import (
	"path/filepath"
)

func TrimExtension(value string) string {
	return value[:len(value)-len(filepath.Ext(value))]
}
