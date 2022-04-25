package elevengo

import (
	"fmt"
	"log"
	"os/exec"
)

func ExampleAgent_CredentialImport() {
	var err error
	agent := Default()

	// Import credential to agent
	if err = agent.CredentialImport(&Credential{
		UID:  "UID-From-Cookie",
		CID:  "CID-From-Cookie",
		SEID: "SEID-From-Cookie",
	}); err != nil {
		log.Fatalf("Import credentail error: %s", err)
	}
	user := agent.User()
	log.Printf("Username: %s", user.Name)
}

func ExampleAgent_FileList() {

	agent := Default()
	if err := agent.CredentialImport(&Credential{
		UID: "", CID: "", SEID: "",
	}); err != nil {
		log.Fatalf("Import credentail error: %s", err)
	}

	cursor, files := &FileCursor{}, make([]*File, 10)
	for cursor.HasMore() {
		n, err := agent.FileList("0", cursor, files)
		if err != nil {
			log.Fatalf("List file failed: %s", err.Error())
		}
		for i := 0; i < n; i++ {
			log.Printf("File: %#v", files[i])
		}
	}
}

func ExampleAgent_Import() {
	var err error

	agent := Default()
	if err = agent.CredentialImport(&Credential{
		UID: "", CID: "", SEID: "",
	}); err != nil {
		log.Fatalf("Import credential failed: %s", err.Error())
	}

	ticket := &ImportTicket{}
	if err = ticket.FromFile("/path/to/local-file"); err != nil {
		log.Fatalf("Init import ticket failed: %s", err.Error())
	}
	if err = agent.Import("0", ticket); err != nil {
		log.Fatalf("Import file to cloud failed: %s", err.Error())
	}
}

func ExampleAgent_OfflineList() {
	//agent := Default()
	// TODO: Import your credentials here

	// Get offline tasks
	//for cursor := OfflineCursor(); cursor.HasMore(); cursor.Next() {
	//	tasks, err := agent.OfflineList(cursor)
	//	if err != nil {
	//		log.Fatalf("Get offline task list error: %s", err)
	//	} else {
	//		for _, task := range tasks {
	//			log.Printf("Offline task: %#v", task)
	//		}
	//	}
	//}
}

func ExampleAgent_DownloadCreateTicket() {
	agent := Default()
	// TODO: Import your credentials here

	// Create download ticket
	ticket := DownloadTicket{}
	err := agent.DownloadCreateTicket("pickcode", &ticket)
	if err != nil {
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
	//agent := Default()
	//// TODO: Import your credentials here
	//
	//filename := "/path/to/file"
	//// Get file info
	//info, err := os.Stat(filename)
	//if err != nil {
	//	log.Fatalf("Get file info error: %s", err)
	//}
	//// Create upload ticket
	//ticket, err := agent.UploadCreateTicket("0", info)
	//if err != nil {
	//	log.Fatalf("Create upload ticket error: %s", err)
	//}
	//// Create temp file to receive upload response
	//tmpFile, err := ioutil.TempFile(os.TempDir(), "115-upload-")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer func() {
	//	_ = os.Remove(tmpFile.Name())
	//}()
	//
	//// Process upload ticket through curl
	//cmd := exec.Command("/usr/bin/curl", ticket.Endpoint, "-o", tmpFile.Name())
	//for name, value := range ticket.Values {
	//	cmd.Args = append(cmd.Args, "-F", fmt.Sprintf("%s=%s", name, value))
	//}
	//// Show upload progress
	//cmd.Args = append(cmd.Args, "-#")
	//// NOTICE: File field should be the LAST one.
	//cmd.Args = append(cmd.Args, "-F", fmt.Sprintf("%s=@%s", ticket.FileField, filename))
	//// Run the command
	//if err = cmd.Run(); err != nil {
	//	log.Fatalf("Execute curl command error: %s", err)
	//}
	//
	//// Parse upload response
	//response, _ := ioutil.ReadAll(tmpFile)
	//file, err := agent.UploadParseResult(response)
	//if err != nil {
	//	log.Fatalf("Parse upload result error: %s", err)
	//} else {
	//	log.Printf("Uploaded file: %#v", file)
	//}
}

func ExampleAgent_VideoGetInfo() {
	agent := Default()
	// TODO: Import your credentials here

	// Get video information
	info := VideoInfo{}
	err := agent.VideoGetInfo("pickcode", &info)
	if err != nil {
		log.Fatalf("Get video info failed: %s", err)
	}
	// Get HLS content
	hls, err := agent.Get(info.PlayUrl)
	if err != nil {
		log.Fatalf("Get HLS content failed: %s", err.Error())
	}
	defer func() {
		_ = hls.Close()
	}()
	// Play HLS through mpv
	cmd := exec.Command("/usr/local/bin/mpv", "-")
	cmd.Stdin = hls
	if err = cmd.Run(); err != nil {
		log.Fatalf("Execute mpv error: %s", err)
	}
}

func ExampleAgent_CaptchaStart() {
	agent := Default()
	// TODO: Import your credentials here

	// Start captcha session.
	session, err := agent.CaptchaStart()
	if err != nil {
		log.Fatalf("Start captcha session error: %s", err)
	}
	// TODO:
	//   * Show `session.CodeImage` and `session.KeysImage` to user.
	//   * Ask user to give the captcha code.

	err = agent.CaptchaSubmit(session, "code")
	if err != nil {
		log.Fatalf("Submit captcha code error: %s", err)
	}
}

func ExampleAgent_QrcodeStart() {
	agent := Default()

	session, err := agent.QrcodeStart()
	if err != nil {
		log.Fatalf("Start QRcode session error: %s", err)
	}
	// TODO:
	// 	Convert `session.Content` to QRcode, show it to user,
	// 	and prompt user to scan it through mobile app.

	for {
		// Get QRcode status
		status, err := agent.QrcodeStatus(session)
		if err != nil {
			if IsQrcodeExpire(err) {
				log.Printf("QRCode expired, please re-generate one.")
				break
			} else {
				log.Fatalf("Get QRcode status error: %s", err)
			}
		} else {
			// Check QRcode status
			if status.IsWaiting() {
				log.Println("Please scan the QRcode in mobile app.")
			} else if status.IsScanned() {
				log.Println("QRcode has beed scanned, please allow this login in mobile app.")
			} else if status.IsAllowed() {
				err = agent.QrcodeLogin(session)
				if err == nil {
					log.Println("QRcode login successed!")
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
