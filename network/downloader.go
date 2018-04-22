package network

import (
	"crypto/tls"
	"net"
	"path/filepath"

	"github.com/fristonio/gvm/logger"
	"github.com/fristonio/gvm/utils"
)

const (
	ACCEPT_RANGE_HEADER   = "Accept-Ranges"
	CONTENT_LENGTH_HEADER = "Content-Length"
)

var log *logger.Logger = logger.New(os.Stdout)

// PartFile Structure
type PartFile struct {
	Url string
	Path string
	RangeFrom int64
	RangeTo int64
}

// Downloader structure - For downloading a file this is the structure
// we need to maintain
type HttpDownloader struct {
	downloadUrl   string
	fileName      string
	sizeDescrip   string
	parts         int64
	contentLength int64
	skipTls       bool
	fileParts     []PartFile
}

var (
	// To get a control over client TLS
	transport = &http.Transport{
		MaxIdleConns:    10,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	// New HTTP client with TLS capabilities
	client = &http.Client{Transport: transport}
)

// Initializes a downloader structure defining a download with values
// and returns it
func NewDownloader(url string, parts int64, skipTls bool) *HttpDownloader {
	log.Infof("New URL for downloading : %s", url)
	req, err := http.NewRequest("GET", url, nil)
	FatalCheck(err)

	res, err := client.Do(req)
	FatalCheck(err)

	if res.Header.Get(ACCEPT_RANGE_HEADER) == "" {
		log.Info("Download url does not support part download, fallback to normal download")
		// Fallback to no part downloading
		parts = 1
	}

	//get download range
	contentLength := res.Header.Get(CONTENT_LENGTH_HEADER)
	if contentLength == "" {
		log.Info("No Content-Length header recieved, fallback to normal download")
		contentLength = "1" //set 1 because of progress bar not accept 0 length
		parts = 1
	}

	log.Info("Starting download with %d connections", parts)
	contentLength, err := strconv.ParseInt(contentLength, 10, 64)
	FatalCheck(err, "Content-Length Header value %s not valid", contentLength)

	sizeDescrip := utils.MemoryBytesToString(contentLength)
	log.Infof("Download Size : %s", sizeDescrip)

	fileName := filepath.Base(url)
	// Final downloader structure
	downloader := HttpDownloader{
		downloadUrl: url,
		fileName: fileName,
		sizeDescrip: sizeDescrip,
		parts: parts,
		contentLength: contentLength,
		skipTls: skipTls,
		fileParts: calculateDownloadParts(int64(parts), contentLength, url)
	}

	return downloader
}

// Takes in the bytes to download and the no of parts and returns and array of Partial File
// Structure which defines each part to be donloaded
func calculateDownloadParts (parts int64, contentLength int64, url string) []PartFile {
	fileParts := make([]FilePart, 0)
	for j := int64(0); j < parts; j++ {
		from := (contentLength / parts) * j
		var to int64
		if j < parts-1 {
			to = (contentLength/parts)*(j+1) - 1
		} else {
			to = contentLength
		}

		file := filepath.Base(url)
		folder := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_DOWNLOAD_DIR)
		if err := MkdirIfNotExist(folder); err != nil {
			log.Fatalf(err)
		}

		fname := fmt.Sprintf("%s.part%d", file, j)
		// ~/.gvm/downloads/fname.part
		path := filepath.Join(folder, fname)
		fileParts = append(fileParts, PartFile{Url: url, Path: path, RangeFrom: from, RangeTo: to})
	}
	return fileParts
}
