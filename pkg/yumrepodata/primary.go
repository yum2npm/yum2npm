package yumrepodata

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"io"
	"regexp"

	"gitlab.com/yum2npm/yum2npm/pkg/utils"
)

type PrimaryMetadata struct {
	XMLName  xml.Name  `xml:"metadata" json:"-"`
	Packages []Package `xml:"package" json:"packages"`
}

type Package struct {
	XMLName     xml.Name       `xml:"package" json:"-"`
	Type        string         `xml:"type,attr" json:"type"`
	Name        string         `xml:"name" json:"name"`
	Arch        string         `xml:"arch" json:"arch"`
	Version     PackageVersion `xml:"version" json:"version"`
	Description string         `xml:"description" json:"description"`
	Summary     string         `xml:"summary" json:"summary"`
}

type PackageVersion struct {
	XMLName xml.Name `xml:"version" json:"-"`
	Epoch   string   `xml:"epoch,attr" json:"epoch"`
	Ver     string   `xml:"ver,attr" json:"ver"`
	Rel     string   `xml:"rel,attr" json:"rel"`
}

func GetPrimary(baseUrl string, repomd RepoMetadata) (PrimaryMetadata, error) {
	var primary string

	for _, x := range repomd.Data {
		if x.Type == "primary" {
			primary = x.Location.Href
			break
		}
	}

	raw, err := utils.FetchUrl(baseUrl + "/" + primary)
	if err != nil {
		return PrimaryMetadata{}, err
	}

	r, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		return PrimaryMetadata{}, err
	}

	res, err := io.ReadAll(r)
	if err != nil {
		return PrimaryMetadata{}, err
	}

	var data PrimaryMetadata

	if err = xml.Unmarshal(res, &data); err != nil {
		return PrimaryMetadata{}, err
	}

	return data, nil
}

func (p PrimaryMetadata) GetPackageVersions(name string) (map[string]ModulePackage, error) {
	var regex = regexp.MustCompile(`.*module.*`)
	filtered := map[string]ModulePackage{}
	for _, i := range p.Packages {
		module := regex.MatchString(i.Version.Rel)
		if i.Name == name && !module {
			v := i.Version.Ver + "-" + i.Version.Rel
			p := ModulePackage{
				Name:    i.Name,
				Version: v,
			}
			filtered[v] = p
		}
	}

	if len(filtered) == 0 {
		return filtered, errors.New("package not found")
	}

	return filtered, nil
}
