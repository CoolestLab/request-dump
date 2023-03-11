package build

import (
	"fmt"
	"runtime"
)

var (
	Version   = "Unknown"
	GoVersion = "Unknown"
	GitHash   = "Unknown"
	BuildTime = "Unknown"
	OSArch    = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	InfoMap   map[string]string
)

func init() {
	InfoMap = map[string]string{
		"version": Version,
		"go":      GoVersion,
		"os/arch": OSArch,
		"commit":  GitHash,
		"built":   BuildTime,
	}
}
