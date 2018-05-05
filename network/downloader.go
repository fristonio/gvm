package network

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fristonio/gvm/logger"
	"github.com/fristonio/gvm/utils"
)

var log *logger.Logger = logger.New(os.Stdout)

// Download the contents of the given URL
func Download(url string, skiptls bool, conn int64, forceClean bool) error {
	var err error
	// We are taking maximum no of concurrent downloads to be conn.

	var files []string

	// Create a singal channel to catch system interrupts
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var isInterrupted = false

	doneChan := make(chan bool, conn)
	fileChan := make(chan string, conn)
	errorChan := make(chan error, 1)
	interruptChan := make(chan bool, conn)

	var downloader *HttpDownloader
	downloader = NewDownloader(url, conn, true)
	// Verfiy and clean already downloaded files in downloads directory.
	if err := downloader.VerifyDownloadDestination(); err != nil {
		log.Errorf("An error occured while verifying download destination : %v", err)
		if forceClean {
			if e := downloader.ClearPreviousDownload(); e != nil {
				return e
			}
		} else {
			return err
		}
	}

	// Start a goroutine for the download
	go downloader.Do(doneChan, fileChan, errorChan, interruptChan)

	for {
		select {
		case <-signal_chan:
			// send parts number of interrupt for each routine
			isInterrupted = true
			for i := int64(0); i < conn; i++ {
				interruptChan <- true
			}
		case file := <-fileChan:
			files = append(files, file)
		case err := <-errorChan:
			log.Errorf("%v", err)
			panic(err) //maybe need better style
		case <-doneChan:
			// Check if the download was successful or it closed due to some  interrupt
			if isInterrupted {
				// Download not finished, interrupt occured. Catch it here
				// As of now we clear the partial downloads when an interrupt occurs
				// But we can give a mechanism to resume the download by saving the state of
				// of current partial downloads.
				log.Warn("Download was interrupted ....")
				log.Warn("Cleaning things up.")
				err = utils.RemoveFilePartials(url)
				utils.FatalCheck(err, "Error occured while removing partial downloads")
				return nil
			} else {
				// Download finished successfully, now join the partial downloads to a single file
				log.Info("Download finished, working on joining partials...")
				err = utils.JoinFilePartials(files, filepath.Base(url))
				utils.FatalCheck(err, "Partial Join of files failed")
				utils.RemoveFilePartials(url)
				utils.FatalCheck(err, "Exitting....")
				return nil
			}
		}
	}
	return nil
}
