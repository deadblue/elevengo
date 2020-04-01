# ElevenGo ![](https://img.shields.io/badge/status-WIP-green.svg)

An API wrapper for 115 NetDisk Service in Golang.

# Exmaple

```
import "github.com/deadblue/elevengo"

// Create API client in default config
client := elevengo.Default()

// Import cookie credentials
cr := &elevengo.Credentials{
	UID: "",
	CID: "",
	SEID: "",
}
if err := client.ImportCredentials(cr); err != nil {
	panic(err)
}

// Get files
files, remain, err := client.FileList("0", 0, 100)
if err != nil {
    panic(err)
}

// Get offline tasks
tasks, remain, err := client.OfflineList(1)
if err != nil {
    panic(err)
}
```

# Features

* Login
  * [x] Import credentials from cookies
  * [x] Login via QRCode
  * [ ] ~~Login via Account/Password~~ (No idea)
* File API
  * [x] List
  * [ ] Search
  * [x] Rename
  * [x] Move
  * [x] Copy
  * [x] Delete
  * [x] Create folder
  * [ ] Download
  * [ ] Upload
* Offline API
  * [x] List tasks
  * [x] Create URL task(s)
  * [x] Delete tasks
  * [x] Clear tasks
* Other
  * [X] Captcha

# License

MIT