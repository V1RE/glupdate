package updaters

import (
	"encoding/json"
	"os/exec"
)

type BrewPackage struct {
	Name              string   `json:"name"`
	InstalledVersions []string `json:"installed_versions"`
	CurrentVersion    string   `json:"current_version"`
}

type BrewPackages struct {
	Formulae []BrewPackage `json:"formulae"`
	Casks    []BrewPackage `json:"casks"`
}

type BrewClient struct {
	path string
}

// List implements Updater.
func (brew *BrewClient) List() ([]Package, error) {
	output, _ := exec.Command(brew.path, "outdated", "--json").Output()

	var packages BrewPackages

	err := json.Unmarshal(output, &packages)
	if err != nil {
		return nil, err
	}

	allPkgs := append(packages.Formulae, packages.Casks...)

	pkgs := make([]Package, 0, len(allPkgs))
	for _, pkg := range allPkgs {
		current := pkg.InstalledVersions[len(pkg.InstalledVersions)-1]
		pkgs = append(pkgs, Package{Name: pkg.Name, Current: current, Latest: pkg.CurrentVersion, Updater: brew})
	}

	return pkgs, nil
}

// Name implements Updater.
func (*BrewClient) Name() string {
	return BrewName
}

const BrewName = "brew"

func Brew() UpdaterFactory {
	return func() (Updater, error) {
		path, err := exec.LookPath(BrewName)
		if err != nil {
			return nil, ErrNotInstalled
		}

		return &BrewClient{path: path}, nil
	}
}
