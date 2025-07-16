package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Package struct {
	ImportPath string
	Dir        string
	GoFiles    []string
}

func ListPackages(patterns []string) ([]Package, error) {
	if len(patterns) == 0 {
		patterns = []string{"./..."}
	}

	args := append([]string{"list", "-json"}, patterns...)
	cmd := exec.Command("go", args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go list failed: %w", err)
	}

	var packages []Package
	decoder := json.NewDecoder(&stdout)

	for decoder.More() {
		var pkg Package
		if err := decoder.Decode(&pkg); err != nil {
			return nil, fmt.Errorf("failed to decode package: %w", err)
		}
		packages = append(packages, pkg)
	}

	return packages, nil
}

func GetGoFilesFromPackages(packages []Package) []string {
	var files []string

	for _, pkg := range packages {
		for _, goFile := range pkg.GoFiles {
			fullPath := fmt.Sprintf("%s/%s", pkg.Dir, goFile)
			files = append(files, fullPath)
		}
	}

	return files
}
