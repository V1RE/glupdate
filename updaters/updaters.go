package updaters

type Package struct {
	Updater Updater
	Name    string
	Current string
	Latest  string
}

type Updater interface {
	List() ([]Package, error)

	Name() string
}

type UpdaterFactory func() (Updater, error)
