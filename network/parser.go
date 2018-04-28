package network

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/fristonio/gvm/utils"
)

type Release struct {
	name        string
	downloadUrl string
}

const (
	TAGS_URL          = "https://go.googlesource.com/go/+refs"
	BASE_DOWNLOAD_URL = "https://go.googlesource.com/go/+archive/%s.tar.gz"
)

// Parses the available release of golang to install
func ParseGoReleases(shouldLog bool) ([]Release, error) {
	log.Info("Releases of go available for download are ")
	releases := make([]Release, 0)

	res, err := http.Get(TAGS_URL)
	defer res.Body.Close()
	if err != nil {
		return releases, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return releases, err
	}
	doc.Find(fmt.Sprintf(".%s", "RefList-item")).Each(func(i int, s *goquery.Selection) {
		releaseName := s.Find("a").Text()

		if utils.GOS_REGEXP.FindString(releaseName) != "" {
			if shouldLog {
				fmt.Println("    " + releaseName)
			}
			releases = append(releases, Release{
				name:        releaseName,
				downloadUrl: fmt.Sprintf(BASE_DOWNLOAD_URL, releaseName),
			})
		}
	})

	return releases, nil
}
