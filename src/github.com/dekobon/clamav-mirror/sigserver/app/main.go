package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

import (
	"github.com/dekobon/clamav-mirror/sigserver"
	"github.com/dekobon/clamav-mirror/utils"
	"github.com/pborman/getopt"
)

var githash = "unknown"
var buildstamp = "unknown"
var appversion = "unknown"

// Main entry point to the server application. This will allow you to run
// the server as a stand-alone binary.
func main() {
	verboseMode, dataFilePath, downloadMirrorURL,
		diffCountThreshold, port, refreshHourInterval := parseCliFlags()

	err := sigserver.RunUpdaterAndServer(verboseMode, dataFilePath, downloadMirrorURL,
		diffCountThreshold, port, refreshHourInterval)

	if err != nil {
		log.Fatal(err)
	}
}

// Function that parses the CLI options passed to the application.
func parseCliFlags() (bool, string, string, uint16, uint16, uint16) {
	verbosePart := getopt.BoolLong("verbose", 'v',
		"Enable verbose mode with additional debugging information")
	versionPart := getopt.BoolLong("version", 'V',
		"Display the version and exit")
	dataFilePart := getopt.StringLong("data-file-path", 'd',
		"/var/clamav/data", "Path to ClamAV data files")
	diffThresholdPart := getopt.Uint16Long("diff-count-threshold", 't',
		100, "Number of diffs to download until we redownload the signature files")
	downloadMirrorPart := getopt.StringLong("download-mirror-url", 'm',
		"http://database.clamav.net", "URL to download signature updates from")
	listenPortPart := getopt.Uint16Long("port", 'p',
		8080, "Port to serve signatures on")
	updateHourlyIntervalPart := getopt.Uint16Long("houry-update-interval", 'h',
		2, "Number of hours to wait between signature updates")

	getopt.Parse()

	if *versionPart {
		fmt.Println("sigupdate")
		fmt.Println("")
		fmt.Printf("Version        : %v\n", appversion)
		fmt.Printf("Git Commit Hash: %v\n", githash)
		fmt.Printf("UTC Build Time : %v\n", buildstamp)
		fmt.Printf("License        : MPLv2\n")

		os.Exit(0)
	}

	if !utils.Exists(*dataFilePart) {
		msg := fmt.Sprintf("Data file path doesn't exist or isn't accessible: %v",
			*dataFilePart)
		log.Fatal(msg)
	}

	dataFileAbsPath, err := filepath.Abs(*dataFilePart)

	if err != nil {
		msg := fmt.Sprintf("Unable to parse absolute path of data file path: %v",
			*dataFilePart)
		log.Fatal(msg)
	}

	if !utils.IsWritable(dataFileAbsPath) {
		msg := fmt.Sprintf("Data file path doesn't have write access for "+
			"current user at path: %v", dataFileAbsPath)
		log.Fatal(msg)
	}

	return *verbosePart, dataFileAbsPath, *downloadMirrorPart, *diffThresholdPart,
		*listenPortPart, *updateHourlyIntervalPart
}
