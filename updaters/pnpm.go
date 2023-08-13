package updaters

import (
	"encoding/json"
	"os/exec"
)

type PnpmPackage struct {
	Current        string `json:"current"`
	Latest         string `json:"latest"`
	Wanted         string `json:"wanted"`
	IsDeprecated   bool   `json:"isDeprecated"`
	DependencyType string `json:"dependencyType"`
}

type PnpmClient struct {
	path string
}

const PnpmName = "pnpm"

// List implements Updater.
func (pnpm *PnpmClient) List() ([]Package, error) {
	output, _ := exec.Command(pnpm.path, "outdated", "--global", "--json").Output()

	var packages map[string]PnpmPackage

	err := json.Unmarshal(output, &packages)
	if err != nil {
		return nil, err
	}

	pkgs := make([]Package, 0, len(packages))
	for name, pp := range packages {
		pkgs = append(pkgs, Package{Name: name, Current: pp.Current, Latest: pp.Latest, Updater: pnpm})
	}

	return pkgs, nil
}

// Name implements Updater.
func (*PnpmClient) Name() string {
	return PnpmName
}

func Pnpm() UpdaterFactory {
	return func() (Updater, error) {
		path, err := exec.LookPath(PnpmName)
		if err != nil {
			return nil, ErrNotInstalled
		}

		return &PnpmClient{path: path}, nil
	}
}
