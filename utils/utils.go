package utils

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sort"

	"github.com/fristonio/gvm/logger"
)

var Log *logger.Logger = logger.New(os.Stdout)

var (
	GVM_ROOT_DIR     string = filepath.Join(os.Getenv("HOME"), ".gvm")
	GVM_DOWNLOAD_DIR string = "downloads"
)

// Returns a string of IPv4 address from a list of IPs returned after lookup
// of a hostname for IPs
func GetIPv4StringArray(ips []net.IP) []string {
	var ipv4s = make([]string, 0)
	for _, ip := range ips {
		// Check if the IP is IPv4 and add it to returned array if it is
		if ip.To4() != nil {
			ipv4s = append(ipv4s, ip.String())
		}
	}
	return ipv4s
}

// Takes byte count in integer format as input and returns a string describing download
// size denoted by the bytecount
func MemoryBytesToString(byteCount int64) string {
	Log.Infof("Bytes : %v", byteCount)
	var downloadSize string
	if byteCount < 1024 {
		downloadSize = fmt.Sprintf("%d Bytes", byteCount)
	} else if byteCount < 1024*1024 {
		downloadSize = fmt.Sprintf("%.1f KBs", float64(byteCount)/1024)
	} else if byteCount < 1024*1024*1024 {
		downloadSize = fmt.Sprintf("%.1f MBs", float64(byteCount)/(1024*1024))
	} else {
		downloadSize = fmt.Sprintf("%.1f GBs", float64(byteCount)/(1024*1024*1024))
	}

	return downloadSize
}

// Takes the name of the folder and create directory as according with
// permissions 755
func MkdirIfNotExist(folder string) error {
	if _, err := os.Stat(folder); err != nil {
		if err = os.MkdirAll(folder, 0755); err != nil {
			return err
		}
	}
	return nil
}

// Remove downloaded file partials corresponding to the url
func RemoveFilePartials(url string) error {
	file := filepath.Base(url)
	downloadsDirectory := filepath.Join(GVM_ROOT_DIR, GVM_DOWNLOAD_DIR)
	files, _ := filepath.Glob(downloadsDirectory + fmt.Sprintf("%s.part*", file))
	err := RemoveAll(files)
	return err
}

func RemoveAll(files []string) error {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// Gets a list of files and joins them into a single file with name out
func JoinFilePartials(files []string, out string) error {
	// Sort the file names so that they are joined in the correct order
	sort.Strings(files)
	downloadsDirectory := filepath.Join(GVM_ROOT_DIR, GVM_DOWNLOAD_DIR)

	Log.Info("Starting to Join file partials")

	outf, err := os.OpenFile(filepath.Join(downloadsDirectory, out), os.O_CREATE|os.O_WRONLY, 0600)
	defer outf.Close()
	if err != nil {
		return err
	}

	for _, f := range files {
		if err = copy(f, outf); err != nil {
			return err
		}
	}
	return nil
}

func copy(from string, to io.Writer) error {
	f, err := os.OpenFile(from, os.O_RDONLY, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	io.Copy(to, f)
	return nil
}
