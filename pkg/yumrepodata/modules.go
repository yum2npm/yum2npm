package yumrepodata

import (
	"context"
	"errors"
	"io"

	"gitlab.com/yum2npm/yum2npm/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Module struct {
	Data ModuleData `yaml:"data"`
}

type ModuleData struct {
	Name        string          `yaml:"name"`
	Stream      string          `yaml:"stream"`
	Description string          `yaml:"description"`
	Artifacts   ModuleArtifacts `yaml:"artifacts"`
}

type ModuleArtifacts struct {
	RPMs []string `yaml:"rpms"`
}

type ModulePackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func GetModules(ctx context.Context, baseUrl string, repomd *RepoMetadata) (mods map[string]map[string]Module, err error) {
	var modules string
	for _, x := range repomd.Data {
		if x.Type == "modules" {
			modules = x.Location.Href
			break
		}
	}

	if len(modules) == 0 {
		return
	}

	var r io.Reader
	r, err = utils.FetchUrl(ctx, baseUrl+"/"+modules)
	if err != nil {
		return
	}

	mods = make(map[string]map[string]Module)

	dec := yaml.NewDecoder(r)
	for {
		var mod Module
		if err := dec.Decode(&mod); err != nil {
			break
		}

		if len(mod.Data.Name) > 0 && len(mod.Data.Stream) > 0 {
			if mods[mod.Data.Name] == nil {
				mods[mod.Data.Name] = make(map[string]Module)
			}
			mods[mod.Data.Name][mod.Data.Stream] = mod
		}
	}

	return mods, nil
}

func (mod Module) GetPackages() []ModulePackage {
	var packagesWithoutArch []string

	for _, p := range mod.Data.Artifacts.RPMs {
		packagesWithoutArch = append(packagesWithoutArch, utils.TrimExtension(p))
	}

	var packageSlice []ModulePackage

	for _, p := range utils.RemoveDuplicateStr(packagesWithoutArch) {
		packageMap := utils.NamedMatches(`^(?P<name>[^\s]+)-(?P<epoch>\d+):(?P<version>[^\s-]+-[^\s-]+)$`, p)
		packageSlice = append(packageSlice, ModulePackage{
			Name:    packageMap["name"],
			Version: packageMap["version"],
		})
	}

	return packageSlice
}

func (mod Module) GetPackageVersions(name string) (map[string]ModulePackage, error) {
	filtered := map[string]ModulePackage{}
	for _, i := range mod.GetPackages() {
		if i.Name == name {
			filtered[i.Version] = i
		}
	}
	if len(filtered) == 0 {
		return filtered, errors.New("package not found")
	}

	return filtered, nil
}
