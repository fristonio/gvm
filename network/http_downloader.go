package network

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/fristonio/gvm/utils"
)

const (
	ACCEPT_RANGE_HEADER   = "Accept-Ranges"
	CONTENT_LENGTH_HEADER = "Content-Length"
)

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
	utils.FatalCheck(err, "Error while making HEAD request to source url")

	res, err := client.Do(req)
	utils.FatalCheck(err, "Error while requesting the resource...")

	if res.Header.Get(ACCEPT_RANGE_HEADER) == "" {
		log.Info("Download url does not support partial download, fallback to normal download")
		// Fallback to no part downloading
		parts = 1
	}

	//get download range
	clen := res.Header.Get(CONTENT_LENGTH_HEADER)
	if clen == "" {
		log.Info("No Content-Length header recieved, fallback to normal download")
		clen = "1"
		parts = 1
	}

	log.Infof("Starting download with %v connections", parts)
	contentLength, err := strconv.ParseInt(clen, 10, 64)
	utils.FatalCheck(err, "Content-Length Header value %s not valid", contentLength)

	sizeDescrip := utils.MemoryBytesToString(contentLength)
	log.Infof("Download Size : %s", sizeDescrip)

	fileName := filepath.Base(url)
	// Final downloader structure
	downloader := &HttpDownloader{
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
	fileParts := make([]PartFile, 0)
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
		if err := utils.MkdirIfNotExist(folder); err != nil {
			log.Fatalf("%v", err)
		}

		fname := fmt.Sprintf("%s.part%d", file, j)
		// ~/.gvm/downloads/fname.part
		path := filepath.Join(folder, fname)
		fileParts = append(fileParts, PartFile{Url: url, Path: path, RangeFrom: from, RangeTo: to})
	}
	return fileParts
}

// Check if the parts and the file does not already exist in the download directory
// Return error if they are already present.
func (d *HttpDownloader) VerifyDownloadDestination() error {
	goSourcePath := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_DOWNLOAD_DIR, d.fileName)
	if _, err := os.Stat(goSourcePath); os.IsNotExist(err) {
		return err
	}

	for _, part := range d.fileParts {
		_, err := os.Stat(part.Path)
		if os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

// Clear/Remove already downloaded parts or file form downloads directory
func (d *HttpDownloader) ClearPreviousDownload() error {
	if err := utils.RemoveFilePartials(d.downloadUrl); err != nil {
		return err
	}
	goSourcePath := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_DOWNLOAD_DIR, d.fileName)
	if err := utils.RemoveAll([1]string{goSourcePath}); err != nil {
		return err
	}
	return nil
}

func (d *HttpDownloader) Do(doneChan chan bool, fileChan chan string, errorChan chan error, interruptChan chan bool) {
	// Sync is for syncronization when implementing concurrency patterns
	// WaitGroup wait for a collection of goroutines to finish
	// The main goroutine calls Add to set the number of goroutines to wait for.
	//Then each of the goroutines runs and calls Done when finished.
	// At the same time, Wait can be used to block until all goroutines have finished.
	var ws sync.WaitGroup

	for i, p := range d.fileParts {
		ws.Add(1)
		// GoRoutine for adding the parts to download
		go func(d *HttpDownloader, partIndex int64, part PartFile) {
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
			req.Header.Add("Range", ranges)

			// Make the above created request
			res, err := client.Do(req)
			if err != nil {
				errorChan <- err
				return
			}
			defer res.Body.Close()
			// Write the contents of the downloads to the file
			f, err := os.OpenFile(part.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)

			defer f.Close()
			if err != nil {
				log.Errorf("%v", err)
				errorChan <- err
				return
			}

			// Make copy interruptable by copy 100 bytes each loop
			current := int64(0)
			var writer io.Writer
			writer = io.MultiWriter(f)
			for {
				select {
				case <-interruptChan:
					return
				default:
					written, err := io.CopyN(writer, res.Body, 100)
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
