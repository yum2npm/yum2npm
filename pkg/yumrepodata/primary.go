package yumrepodata

import (
	"context"
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

func GetPrimary(ctx context.Context, baseUrl string, repomd *RepoMetadata) (meta *PrimaryMetadata, err error) {
	var primary string

	for _, x := range repomd.Data {
		if x.Type == "primary" {
			primary = x.Location.Href
			break
		}
	}

	if len(primary) == 0 {
		return
	}

	var r io.Reader
	r, err = utils.FetchUrl(ctx, baseUrl+"/"+primary)
	if err != nil {
		return
	}

	meta = &PrimaryMetadata{}

	err = xml.NewDecoder(r).Decode(meta)

	return
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
