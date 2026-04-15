package metadata

import (
	"runtime/debug"
	"sync"
)

var (
	once    sync.Once
	version = "devel"
)

// GetVersion returns the cached library version.
func GetVersion() string {
	once.Do(func() {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}

		const myModulePath = "github.com/aleal/zero"

		// 1. Check if we are a dependency of another module
		for _, dep := range info.Deps {
			if dep.Path == myModulePath {
				version = dep.Version
				return
			}
		}

		// 2. Check if the main binary is this module itself (during dev)
		if info.Main.Path == myModulePath && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	})

	return version
}
