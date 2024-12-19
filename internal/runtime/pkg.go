package runtime

type Package interface {
	Run(fn string, args []variable) ([]funcReturn, error)
}

func NewPackage(name string) Package {
	switch name {
	case "io":
		return IO{}
	case "strs":
		return Strs{}
	}
	return nil
}
