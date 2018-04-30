package utils

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/fristonio/gvm/logger"
)

var Log *logger.Logger = logger.New(os.Stdout)

var GOS_REGEXP *regexp.Regexp = getGosRegexp()

func getGosRegexp() *regexp.Regexp {
	gosRegexp, _ := regexp.Compile(`^go[\d\.]+$`)
	return gosRegexp
}

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

func PrintInstalledGos(gos []string) {
	if len(gos) == 0 {
		Log.Warn("No gos installed, to view a list of versions available use: go list-remote")
		return
	}
	for i, f := range gos {
		fmt.Println(strconv.Itoa(i+1) + ". " + f)
	}
}

// Checks if a directory is present, if it is return no error if not
// Create the provide directory structure provided in dirString creating all necessery parents
func CreateDirStrucutre(dirString string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirString, 0660)
		return err
	}
}

func CheckIfDirExist(dirString string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return err
	}
}
