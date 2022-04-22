package elevengo

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/deadblue/elevengo/internal/protocol"
	"github.com/deadblue/elevengo/internal/webapi"
	"github.com/deadblue/gostream/quietly"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	importPreSize = 128 * 1024
)

type ImportTicket struct {
	PreId   string
	QuickId string
	Name    string
	Size    int64
}

func (t *ImportTicket) From(r io.Reader, name string, size int64) (err error) {
	t.Name = name
	t.Size = size

	// Calculate hash
	hash := sha1.New()

	// Pre ID
	if size > importPreSize {
		_, err = io.CopyN(hash, r, importPreSize)
	} else {
		_, err = io.Copy(hash, r)
	}
	if err != nil {
		return
	}
	t.PreId = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))

	// Quick ID
	if size > importPreSize {
		if _, err = io.Copy(hash, r); err != nil {
			return
		}
		t.QuickId = strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	} else {
		t.QuickId = t.PreId
	}
	return
}

func (t *ImportTicket) FromFile(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer quietly.Close(file)
	info, err := file.Stat()
	if err != nil {
		return
	}
	return t.From(file, info.Name(), info.Size())
}

func (a *Agent) Import(dirId string, ticket *ImportTicket) (err error) {
	if !a.ut.Available() {
		if err = a.uploadInit(); err != nil {
			return
		}
	}
	// Prepare request
	targetId := fmt.Sprintf("U_1_%s", dirId)
	qs := protocol.Params{}.
		With("appid", a.ut.AppId).
		With("appversion", a.ut.AppVer).
		WithInt("isp", a.ut.IspType).
		With("sig", a.updateCalculateSignature(targetId, ticket.QuickId)).
		With("format", "json").
		WithNow("t")
	form := protocol.Params{}.
		With("app_ver", a.ut.AppVer).
		With("preid", ticket.PreId).
		With("quickid", ticket.QuickId).
		With("target", targetId).
		With("fileid", ticket.QuickId).
		With("filename", ticket.Name).
		WithInt64("filesize", ticket.Size).
		WithInt("userid", a.user.Id)
	// Send request
	resp := &webapi.UploadInitResponse{}
	if err = a.wc.CallJsonApi(webapi.ApiUploadInit, qs, form, resp); err != nil {
		return
	}
	return
}

func (a *Agent) updateCalculateSignature(targetId string, fileId string) string {
	buf := &bytes.Buffer{}
	buf.WriteString(strconv.Itoa(a.user.Id))
	buf.WriteString(fileId)
	buf.WriteString(fileId)
	buf.WriteString(targetId)
	buf.WriteRune('0')
	digest := sha1.Sum(buf.Bytes())

	buf.Reset()
	buf.WriteString(a.ut.UserKey)
	buf.WriteString(hex.EncodeToString(digest[:]))
	buf.WriteString("000000")
	digest = sha1.Sum(buf.Bytes())
	return strings.ToUpper(hex.EncodeToString(digest[:]))
}
