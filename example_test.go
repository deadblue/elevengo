package elevengo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/deadblue/elevengo/option"
)

func ExampleAgent_CredentialImport() {
	agent := Default()

	// Import credential to agent
	if err := agent.CredentialImport(&Credential{
		UID:  "UID-From-Cookie",
		CID:  "CID-From-Cookie",
		SEID: "SEID-From-Cookie",
	}); err != nil {
		log.Fatalf("Import credentail error: %s", err)
	}
}

func ExampleAgent_FileIterate() {
	agent := Default()

	it, err := agent.FileIterate("0")
	if err != nil {
		log.Fatalf("Iterate file failed: %s", err.Error())
	}
	log.Printf("File count: %d", it.Count())
	for index, file := range it.Items() {
		log.Printf("File: %d => %#v", index, file)
	}
}

func ExampleAgent_OfflineIterate() {
	agent := Default()
	it, err := agent.OfflineIterate()
	if err == nil {
		log.Printf("Task count: %d", it.Count())
		for index, task := range it.Items() {
			log.Printf("Offline task: %d => %#v", index, task)
		}
	} else {
		log.Fatalf("Iterate offline task failed: %s", err)
	}
}

func ExampleAgent_DownloadCreateTicket() {
	agent := Default()

	// Create download ticket
	var err error
	ticket := DownloadTicket{}
	if err = agent.DownloadCreateTicket("pickcode", &ticket); err != nil {
		log.Fatalf("Get download ticket error: %s", err)
	}

	// Process download ticket through curl
	cmd := exec.Command("/usr/bin/curl", ticket.Url)
	for name, value := range ticket.Headers {
		cmd.Args = append(cmd.Args, "-H", fmt.Sprintf("%s: %s", name, value))
	}
	cmd.Args = append(cmd.Args, "-o", ticket.FileName)
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("File downloaded to %s", ticket.FileName)
	}
}

func ExampleAgent_UploadCreateTicket() {
	agent := Default()

	filename := "/path/to/file"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Open file failed: %s", err.Error())
	}

	ticket := &UploadTicket{}
	if err = agent.UploadCreateTicket("dirId", path.Base(filename), file, ticket); err != nil {
		log.Fatalf("Create upload ticket failed: %s", err.Error())
	}
	if ticket.Exist {
		log.Printf("File already exists!")
		return
	}

	// Make temp file to receive upload result
	tmpFile, err := os.CreateTemp("", "curl-upload-*")
	if err != nil {
		log.Fatalf("Create temp file failed: %s", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	// Use "curl" to upload file
	cmd := exec.Command("curl", ticket.Url,
		"-o", tmpFile.Name(), "-#",
		"-T", filename)
	for name, value := range ticket.Header {
		cmd.Args = append(cmd.Args, "-H", fmt.Sprintf("%s: %s", name, value))
	}
	if err = cmd.Run(); err != nil {
		log.Fatalf("Upload failed: %s", err)
	}

	// Parse upload result
	uploadFile := &File{}
	if err = agent.UploadParseResult(tmpFile, uploadFile); err != nil {
		log.Fatalf("Parse upload result failed: %s", err)
	} else {
		log.Printf("Uploaded file: %#v", file)
	}
}

func ExampleAgent_VideoCreateTicket() {
	agent := Default()

	// Create video play ticket
	ticket := &VideoTicket{}
	err := agent.VideoCreateTicket("pickcode", ticket)
	if err != nil {
		log.Fatalf("Get video info failed: %s", err)
	}

	headers := make([]string, 0, len(ticket.Headers))
	for name, value := range ticket.Headers {
		headers = append(headers, fmt.Sprintf("'%s: %s'", name, value))
	}

	// Play HLS through mpv
	cmd := exec.Command("mpv")
	cmd.Args = append(cmd.Args,
		fmt.Sprintf("--http-header-fields=%s", strings.Join(headers, ",")),
		ticket.Url,
	)
	cmd.Run()
}

func ExampleAgent_QrcodeStart() {
	agent := Default()

	session := &QrcodeSession{}
	err := agent.QrcodeStart(session)
	if err != nil {
		log.Fatalf("Start QRcode session error: %s", err)
	}
	// TODO: Show QRcode and ask user to scan it via 115 app.
	for done := false; !done && err != nil; {
		done, err = agent.QrcodePoll(session)
	}
	if err != nil {
		log.Fatalf("QRcode login failed, error: %s", err)
	}
}

func ExampleAgent_Import() {
	var err error

	// Initialize two agents for sender and receiver
	sender, receiver := Default(), Default()
	sender.CredentialImport(&Credential{
		UID: "", CID: "", SEID: "",
	})
	receiver.CredentialImport(&Credential{
		UID: "", CID: "", SEID: "",
	})

	// File to send on sender's storage
	fileId := "12345678"
	// Create import ticket by sender
	ticket, pickcode := &ImportTicket{}, ""
	if pickcode, err = sender.ImportCreateTicket(fileId, ticket); err != nil {
		log.Fatalf("Get file info failed: %s", err)
	}

	// Directory to save file on receiver's storage
	dirId := "0"
	// Call Import first time
	if err = receiver.Import(dirId, ticket); err != nil {
		if ie, ok := err.(*ErrImportNeedCheck); ok {
			// Calculate sign value by sender
			signValue, err := sender.ImportCalculateSignValue(pickcode, ie.SignRange)
			if err != nil {
				log.Fatalf("Calculate sign value failed: %s", err)
			}
			// Update ticket and import again
			ticket.SignKey, ticket.SignValue = ie.SignKey, signValue
			if err = receiver.Import(dirId, ticket); err == nil {
				log.Print("Import succeeded!")
			} else {
				log.Fatalf("Import failed: %s", err)
			}
		} else {
			log.Fatalf("Import failed: %s", err)
		}
	} else {
		log.Print("Import succeeded!")
	}
}

func ExampleNew() {
	// Customize agent
	agent := New(
		// Custom agent name
		option.AgentNameOption("Evangelion/1.0"),
		// Sleep 100~500 ms between two API calling
		option.AgentCooldownOption{Min: 100, Max: 500},
	)

	var err error
	if err = agent.CredentialImport(&Credential{
		UID: "", CID: "", SEID: "",
	}); err != nil {
		log.Fatalf("Invalid credential, error: %s", err)
	}
}
