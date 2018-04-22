package utils

import (
	"fmt"
	"net"
)

const (
	GVM_ROOT_DIR     string = "~/.gvm"
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
	if byteCount < 1024 {
		downloadSize := fmt.Sprintf("%d Bytes", byteCount)
	} else if byteCount < 1024*1024 {
		downloadSize := fmt.Sprintf("%.1f KBs", float64(byteCount)/1024)
	} else if byteCount < 1024*1024 {
		downloadSize := fmt.Sprintf("%.1f MBs", float64(byteCount)/(1024*1024))
	} else {
		downloadSize := fmt.Sprintf("%.1f GBs", float64(byteCount)/(1024*1024*1024))
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
