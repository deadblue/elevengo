package elevengo

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/internal/webapi"
)

var (
	regexpImportURI = regexp.MustCompile(`115://\|(\w+)\|(\d+)\|(\w+)\|(\w+)\|/`)
)

type ImportTicket struct {
	// File name to store
	Name string
	// File size
	Size int64
	// SHA1 hash of the first 128KB content, in upper-case HEX format
	PreId string
	// SHA1 hash of the whole file, in upper-case HEX format
	QuickId string
}

// FromFile fills ticket from a local file.
func (t *ImportTicket) FromFile(name string) (err error) {
	file, err := os.Open(name)
	if err != nil {
		return
	}
	info, err := file.Stat()
	if err != nil {
		return
	} else if info.IsDir() {
		return webapi.ErrImportDirectory
	}
	defer util.QuietlyClose(file)
	return t.From(path.Base(name), file)
}

// From fills ticket with given name and r.
func (t *ImportTicket) From(name string, r io.Reader) (err error) {
	dr := hash.DigestResult{}
	if err = hash.Digest(r, &dr); err != nil {
		return
	}
	t.Name = name
	t.Size = dr.Size
	t.PreId = dr.PreId
	t.QuickId = dr.QuickId
	return
}

func (t *ImportTicket) FromURI(uri string) error {
	fields := regexpImportURI.FindStringSubmatch(uri)
	if len(fields) == 0 {
		return webapi.ErrInvalidImportURI
	}
	var err error
	if t.Name, err = url.QueryUnescape(fields[1]); err != nil {
		t.Name = ""
		return webapi.ErrInvalidImportURI
	}
	t.Size, _ = strconv.ParseInt(fields[2], 10, 64)
	t.QuickId, t.PreId = fields[3], fields[4]
	return nil
}

func (t *ImportTicket) ToURI() string {
	return fmt.Sprintf("115://|%s|%d|%s|%s|/",
		url.QueryEscape(t.Name), t.Size, t.QuickId, t.PreId)
}

// Import imports file which already exists on cloud to your account.
func (a *Agent) Import(dirId string, ticket *ImportTicket) (err error) {
	var exist bool
	exist, err = a.uploadInit(dirId, ticket.Name, ticket.Size, ticket.PreId, ticket.QuickId, nil)
	if err == nil && !exist {
		err = webapi.ErrNotExist
	}
	return
}

// ImportCreateTicket creates an ImportTicket from fileId.
func (a *Agent) ImportCreateTicket(fileId string, ticket *ImportTicket) (err error) {
	// Get file information
	file := &File{}
	if err = a.FileGet(fileId, file); err != nil {
		return err
	}
	if file.IsDirectory {
		return webapi.ErrImportDirectory
	}
	// Fill ImportTicket
	ticket.Name = file.Name
	ticket.Size = file.Size
	ticket.QuickId = file.Sha1
	if file.Size <= webapi.UploadPreSize {
		ticket.PreId = file.Sha1
		return
	}
	// Get first 128K data of the file
	dt := &DownloadTicket{}
	if err = a.DownloadCreateTicket(file.PickCode, dt); err != nil {
		return
	}
	preBody, err := a.GetRange(dt.Url, webapi.UploadPreSize, 0)
	if err != nil {
		return
	}
	defer util.QuietlyClose(preBody)
	ticket.PreId = hash.Sha1HexUpper(preBody)
	return
}
