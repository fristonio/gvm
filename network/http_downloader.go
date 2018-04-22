package network

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

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
	Url       string
	Path      string
	RangeFrom int64
	RangeTo   int64
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
	req, err := http.NewRequest("HEAD", url, nil)
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
		downloadUrl:   url,
		fileName:      fileName,
		sizeDescrip:   sizeDescrip,
		parts:         parts,
		contentLength: contentLength,
		skipTls:       skipTls,
		fileParts:     calculateDownloadParts(int64(parts), contentLength, url),
	}

	return downloader
}

// Takes in the bytes to download and the no of parts and returns and array of Partial File
// Structure which defines each part to be donloaded
func calculateDownloadParts(parts int64, contentLength int64, url string) []PartFile {
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

func (d *HttpDownloader) Download(doneChan chan bool, fileChan chan string, errorChan chan error, interruptChan chan bool) {
	// Sync is for syncronization when implementing concurrency patterns
	// WaitGroup wait for a collection of goroutines to finish
	// The main goroutine calls Add to set the number of goroutines to wait for.
	//Then each of the goroutines runs and calls Done when finished.
	// At the same time, Wait can be used to block until all goroutines have finished.
	var ws sync.WaitGroup
	var err error

	for i, p := range d.fileParts {
		ws.Add(1)
		// GoRoutine for adding the parts to download
		go func(d *HttpDownloader, partIndex int64, part FilePart) {
			// Call done when the routine execution finish, to let wait group know about it.
			defer ws.Done()

			var ranges string
			// Ranges Header for a part of the download
			if part.RangeTo != d.contentLength {
				ranges = fmt.Sprintf("bytes=%d-%d", part.RangeFrom, part.RangeTo)
			} else {
				ranges = fmt.Sprintf("bytes=%d-", part.RangeFrom) //get all
			}

			// Send the GET request
			req, err := http.NewRequest("GET", d.downloadUrl, nil)
			// If an error occurs push that error to error channel
			if err != nil {
				errorChan <- err
				return
			}

			// Add range header in cases when part downloading is possible
			if d.parts > 1 {
				req.Header.Add("Range", ranges)
				if err != nil {
					errorChan <- err
					return
				}
			}

			// Make the above created request
			res, err := client.Do(req)
			if err != nil {
				errorChan <- err
				return
			}
			defer res.Body.Close()
			// Write the contents of the downloads to the file
			f, err := os.OpenFile(part.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

			defer f.Close()
			if err != nil {
				log.Errorf("%v\n", err)
				errorChan <- err
				return
			}

			// Make copy interruptable by copy 100 bytes each loop
			current := int64(0)
			for {
				select {
				case <-interruptChan:
					return
				default:
					written, err := io.CopyN(writer, resp.Body, 100)
					current += written
					if err != nil {
						if err != io.EOF {
							errorChan <- err
						}
						// File part download completes here
						fileChan <- part.Path
						return
					}
				}
			}
		}(d, int64(i), p)
	}

	// Wait here until all the goroutines are done
	ws.Wait()
	doneChan <- true
}
