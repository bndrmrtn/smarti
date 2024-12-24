package runtime

import "github.com/bndrmrtn/smarti/internal/packages"

func NewPackage(name string) packages.Package {
	switch name {
	case "io":
		return packages.IO{}
	case "strs":
		return packages.Strs{}
	case "numbers":
		return packages.Numbers{}
	case "env":
		return packages.Env{}
	}
	return nil
}
