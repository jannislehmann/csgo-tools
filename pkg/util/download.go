package util

import (
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// ErrInvalidDownloadURL is returned when the url to be download is invalid or malicious.
var ErrInvalidDownloadURL = errors.New("invalid download url")

// DemoNotFoundError is used when a valid matchid / demo is not found or can no longer be downloaded.
type errDemoNotFound struct {
	URL string
}

func (e errDemoNotFound) Error() string {
	const msg = "demo no longer downloadable: %s"
	return fmt.Sprintf(msg, e.URL)
}

func IsDemoNotFoundError(err error) bool {
	_, ok := err.(errDemoNotFound)
	return ok
}

// DownloadDemo will download a demo from an url and decompress and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadDemo(url string, demoDir string, lastModified time.Time) error {
	// Validate the url
	reValve := regexp.MustCompile(`^http:\/\/replay[\d]{3}\.valve\.net\/730\/[\d]{21}_([\d]*)\.dem\.bz2$`)
	reFaceit := regexp.MustCompile(`^https:\/\/demos-([\w]*)-([\w]*)\.faceit-cdn\.net\/csgo\/[\d]{1}-\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b-[\d]{1}-[\d]{1}\.dem\.gz$`)

	if !reValve.MatchString(url) && !reFaceit.MatchString(url) {
		return ErrInvalidDownloadURL
	}

	// Get file name.
	filename := strings.Split(path.Base(url), ".")[0] + ".dem"
	filePath := path.Join(demoDir, filename)

	// Get the data.
	resp, err := http.Get(url) //nolint // We have to take dynamic replay urls in order to download them. URL is validated before.
	if err != nil || resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return errDemoNotFound{URL: url}
	}

	defer resp.Body.Close()

	// Create the file.
	out, err := os.Create(filePath)
	if err != nil {
		out.Close()
		log.Error(err)
		return err
	}
	defer out.Close()

	var cr io.Reader
	if strings.HasSuffix(url, "gz") {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		cr = reader
	} else if strings.HasSuffix(url, "bz2") {
		cr = bzip2.NewReader(resp.Body)
	} else {
		return errors.New("unknown file compression")
	}

	// Decompress and write to file.
	if _, err = io.Copy(out, cr); err != nil {
		return err
	}

	// Close file before trying to update the modified information.
	out.Close()

	// Update file modified information.
	if err = os.Chtimes(filePath, lastModified, lastModified); err != nil {
		const msg = "unable to set correct last modified date for demo %v"
		log.Warnf(msg, filename)
		log.Error(err)
	}

	return nil
}
