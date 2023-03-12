package elevengo

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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
	for ; err == nil; err = it.Next() {
		file := &File{}
		_ = it.Get(file)
		log.Printf("File: %d => %#v", it.Index(), file)
	}
	if !IsIteratorEnd(err) {
		log.Fatalf("Iterate file failed: %s", err.Error())
	}
}

func ExampleAgent_OfflineIterate() {
	agent := Default()

	for it, err := agent.OfflineIterate(); err == nil; err = it.Next() {
		task := &OfflineTask{}
		err = it.Get(task)
		if err != nil {
			log.Printf("Offline task: %#v", task)
		}
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

func ExampleAgent_VideoGet() {
	agent := Default()

	// Get video information
	info := Video{}
	err := agent.VideoGet("pickcode", &info)
	if err != nil {
		log.Fatalf("Get video info failed: %s", err)
	}

	// Get HLS content
	hlsData, err := agent.Get(info.PlayUrl)
	if err != nil {
		log.Fatalf("Get HLS content failed: %s", err.Error())
	}
	defer func() {
		_ = hlsData.Close()
	}()

	// Play HLS through mpv
	cmd := exec.Command("mpv", "-")
	cmd.Stdin = hlsData
	if err = cmd.Run(); err != nil {
		log.Fatalf("Execute mpv error: %s", err)
	}
}

func ExampleAgent_CaptchaStart() {
	agent := Default()

	var err error
	// Start captcha session.
	session := &CaptchaSession{}
	if err = agent.CaptchaStart(session); err != nil {
		log.Fatalf("Start captcha session error: %s", err)
	}

	// 1. Show `session.CodeImage` and `session.KeysImage` to user.
	// 2. Ask user to give the captcha code.

	if err = agent.CaptchaSubmit(session, "code"); err != nil {
		log.Fatalf("Submit captcha code error: %s", err)
	}
}

func ExampleAgent_QrcodeStart() {
	agent := Default()

	session := &QrcodeSession{}
	err := agent.QrcodeStart(session)
	if err != nil {
		log.Fatalf("Start QRcode session error: %s", err)
	}
	// Convert `session.Content` to QRCode, show it to user, and prompt user
	// to scan it using mobile app.

	for {
		var status QrcodeStatus
		// Get QR-Code status
		status, err = agent.QrcodeStatus(session)
		if err != nil {
			log.Fatalf("Get QRCode status error: %s", err)
		} else {
			// Check QRCode status
			if status.IsWaiting() {
				log.Println("Please scan the QRCode in mobile app.")
			} else if status.IsScanned() {
				log.Println("QRCode has been scanned, please allow this login in mobile app.")
			} else if status.IsAllowed() {
				err = agent.QrcodeLogin(session)
				if err == nil {
					log.Println("QRCode login successes!")
				} else {
					log.Printf("Submit QRcode login error: %s", err)
				}
				break
			} else if status.IsCanceled() {
				fmt.Println("User canceled this login!")
				break
			}
		}
	}

}
