package elevengo

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func ExampleAgent_CredentialImport() {
	agent := Default()
	// Import credential to agent
	err := agent.CredentialImport(&Credential{
		UID:  "UID-from-cookie",
		CID:  "CID-from-cookie",
		SEID: "SEID-from-cookie",
	})
	if err != nil {
		log.Fatalf("Import credentail error: %s", err)
	} else {
		user := agent.User()
		log.Printf("Username: %s", user.Name)
	}
}

func ExampleAgent_FileList() {
	agent := Default()
	// TODO: Import your credentials here

	// Get files under root directory
	for cursor := FileCursor(); cursor.HasMore(); cursor.Next() {
		files, err := agent.FileList("0", cursor)
		if err != nil {
			log.Fatalf("Get file list error: %s", err)
		} else {
			for _, file := range files {
				log.Printf("Remote file: %#v", file)
			}
		}
	}
}

func ExampleAgent_OfflineList() {
	agent := Default()
	// TODO: Import your credentials here

	// Get offline tasks
	for cursor := OfflineCursor(); cursor.HasMore(); cursor.Next() {
		tasks, err := agent.OfflineList(cursor)
		if err != nil {
			log.Fatalf("Get offline task list error: %s", err)
		} else {
			for _, task := range tasks {
				log.Printf("Offline task: %#v", task)
			}
		}
	}
}

func ExampleAgent_DownloadCreateTicket() {
	agent := Default()
	// TODO: Import your credentials here

	// Create download ticket
	ticket, err := agent.DownloadCreateTicket("pickcode")
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
	agent := Default()
	// TODO: Import your credentials here

	filename := "/path/to/file"
	// Get file info
	info, err := os.Stat(filename)
	if err != nil {
		log.Fatalf("Get file info error: %s", err)
	}
	// Create upload ticket
	ticket, err := agent.UploadCreateTicket("0", info)
	if err != nil {
		log.Fatalf("Create upload ticket error: %s", err)
	}
	// Createa temp file to receive upload response
	tmpFile, err := ioutil.TempFile(os.TempDir(), "115-upload-")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	// Process upload ticket through curl
	cmd := exec.Command("/usr/bin/curl", ticket.Endpoint, "-o", tmpFile.Name())
	for name, value := range ticket.Values {
		cmd.Args = append(cmd.Args, "-F", fmt.Sprintf("%s=%s", name, value))
	}
	// Show upload progress
	cmd.Args = append(cmd.Args, "-#")
	// NOTICE: File field should be the LAST one.
	cmd.Args = append(cmd.Args, "-F", fmt.Sprintf("%s=@%s", ticket.FileField, filename))
	// Run the command
	if err = cmd.Run(); err != nil {
		log.Fatalf("Execute curl command error: %s", err)
	}

	// Parse upload response
	response, _ := ioutil.ReadAll(tmpFile)
	file, err := agent.UploadParseResult(response)
	if err != nil {
		log.Fatalf("Parse upload result error: %s", err)
	} else {
		log.Printf("Uploaded file: %#v", file)
	}
}
