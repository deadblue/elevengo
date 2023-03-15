package elevengo

import (
	"fmt"

	"github.com/deadblue/elevengo/internal/webapi"
)

type errImportRequireCheck struct {
	signKey   string
	signRange string
}

func (e *errImportRequireCheck) Error() string {
	return ""
}

// ImportTicket container reqiured fields to import(aka. quickly upload) a file
// to your 115 cloud storage.
type ImportTicket struct {
	// File base name
	FileName  string
	// File size in bytes
	FileSize  int64
	// File SHA-1 hash, in upper-case HEX format
	FileSha1  string
	// Sign key from 115 server.
	SignKey string
	// SHA-1 hash value of a part in file, in upper-case HEX format
	SignValue string
}

/*
Import tries to import a file to your 115 cloud storage.

Example:
	agent := Default()
	ticket := &ImportTicket{
		FileName: "hello.mp4",
		FileSize: 12345678,
		FileSha1: "0123456789ABCDEF0123456789ABCDEF01234567"
	}
	err := agent.Import("0", ticket)
	if ok, key, rng := IsImportCheckRequired(err); ok {
		ticket.SignKey = key
		// TODO: Implement CalculateHashRange
		ticket.SignValue = CalculateHashRange(file, rng)
		err = agnet.Import("0", ticket)
	}
	if err != nil {
		log.Fatal
	}

*/
func (a *Agent) Import(dirId string, ticket *ImportTicket) (err error) {
	if err = a.uploadInitHelper(); err != nil {
		return
	}
	target := fmt.Sprintf("U_1_%s", dirId)
	initData := &webapi.UploadInitData{
		FileId: ticket.FileSha1,
		FileName: ticket.FileName,
		FileSize: ticket.FileSize,
		Target: target,
		Signature: a.uh.CalculateSignature(ticket.FileSha1, target),
		SignKey: ticket.SignKey,
		SignValue: ticket.SignValue,
	}
	exist, checkRange := false, ""
	if exist, checkRange, err = a.uploadInitInternal(initData, nil); err == nil {
		if checkRange != "" {
			err = &errImportRequireCheck{
				signKey: initData.SignKey,
				signRange: checkRange,
			}
		} else if !exist {
			err = webapi.ErrNotExist
		}
	}
	return
}

// IsImportCheckRequired
func IsImportCheckRequired(err error) (ok bool, key, rng string) {
	e, ok := err.(*errImportRequireCheck)
	if ok {
		key, rng = e.signKey, e.signRange
	}
	return
}
