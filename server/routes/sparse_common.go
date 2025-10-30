//go:build !windows

package routes

import "os"

func setSparse(*os.File) {
}
