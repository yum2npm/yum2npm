package yumrepodata

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"

	"github.com/h2non/filetype"

	"github.com/ulikunitz/xz"
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

func GetModules(baseUrl string, repomd RepoMetadata) (map[string]map[string]Module, error) {
	var modules string
	for _, x := range repomd.Data {
		if x.Type == "modules" {
			modules = x.Location.Href
			break
		}
	}

	if len(modules) == 0 {
		return make(map[string]map[string]Module), nil
	}

	raw, err := utils.FetchUrl(baseUrl + "/" + modules)
	if err != nil {
		return make(map[string]map[string]Module), err
	}

	var r io.Reader

	kind, err := filetype.Match(raw)
	if err != nil {
		return make(map[string]map[string]Module), err
	}
	switch kind.MIME.Value {
	case "application/gzip":
		r, err = gzip.NewReader(bytes.NewReader(raw))
	case "application/x-xz":
		r, err = xz.NewReader(bytes.NewReader(raw))
	default:
		r = bytes.NewReader(raw)
		err = nil
	}

	if err != nil {
		return make(map[string]map[string]Module), err
	}

	res, err := io.ReadAll(r)
	if err != nil {
		return make(map[string]map[string]Module), err
	}

	data := make(map[string]map[string]Module)
	var mods []Module

	dec := yaml.NewDecoder(bytes.NewReader(res))
	for {
		var mod Module
		if err := dec.Decode(&mod); err != nil {
			break
		}
		mods = append(mods, mod)
	}

	for _, mod := range mods {
		if len(mod.Data.Name) > 0 && len(mod.Data.Stream) > 0 {
			if data[mod.Data.Name] == nil {
				data[mod.Data.Name] = make(map[string]Module)
			}
			data[mod.Data.Name][mod.Data.Stream] = mod
		}
	}

	return data, nil
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
