package elevengo

import (
	"errors"
	"github.com/deadblue/elevengo/internal/crypto/hash"
	"github.com/deadblue/elevengo/internal/util"
	"io"
	"os"
	"path"
)

var (
	ErrImportFailed = errors.New("import failed: file does not exist no remote")
)

type ImportTicket struct {
	Name    string
	Size    int64
	PreId   string
	QuickId string
}

// FromFile fills ticket from a local file.
func (t *ImportTicket) FromFile(name string) (err error) {
	file, err := os.Open(name)
	if err != nil {
		return
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

// Import imports file which already exists on cloud to your account.
func (a *Agent) Import(dirId string, ticket *ImportTicket) (err error) {
	var exist bool
	exist, err = a.uploadInit(dirId, ticket.Name, ticket.Size, ticket.PreId, ticket.QuickId, nil)
	if err == nil && !exist {
		err = ErrImportFailed
	}
	return
}
