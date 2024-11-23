package yumrepodata

import (
	"bytes"
	"context"
	"encoding/xml"
	"gitlab.com/yum2npm/yum2npm/pkg/utils"
	"io"
)

type RepoMetadata struct {
	XMLName  xml.Name     `xml:"repomd"`
	Revision string       `xml:"revision"`
	Tags     RepomdTags   `xml:"tags"`
	Data     []RepomdData `xml:"data"`
}

type RepomdTags struct {
	XMLName xml.Name `xml:"tags"`
	Distro  string   `xml:"distro"`
}

type RepomdData struct {
	XMLName      xml.Name           `xml:"data"`
	Type         string             `xml:"type,attr"`
	Checksum     RepomdChecksum     `xml:"checksum"`
	OpenChecksum RepomdOpenChecksum `xml:"open-checksum"`
	Location     RepomdLocation     `xml:"location"`
	Timestamp    int                `xml:"timestamp"`
	Size         int                `xml:"size"`
	OpenSize     int                `xml:"open-size"`
}

type RepomdChecksum struct {
	XMLName xml.Name `xml:"checksum"`
	Type    string   `xml:"type,attr"`
}

type RepomdOpenChecksum struct {
	XMLName xml.Name `xml:"open-checksum"`
	Type    string   `xml:"type,attr"`
}

type RepomdLocation struct {
	XMLName xml.Name `xml:"location"`
	Href    string   `xml:"href,attr"`
}

func GetRepoMetadata(ctx context.Context, baseUrl string) (repomd *RepoMetadata, err error) {
	r, err := utils.FetchUrl(ctx, baseUrl+"/repodata/repomd.xml")
	if err != nil {
		return
	}

	repomd = &RepoMetadata{}

	b, err := io.ReadAll(r)

	err = xml.NewDecoder(bytes.NewReader(b)).Decode(repomd)

	return
}
