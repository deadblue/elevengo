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
  * [ ] Login via QRCode
  * [ ] Login via Account/Password (No idea now)
* File API
  * [x] List
  * [x] Search
  * [x] Rename
  * [x] Move
  * [x] Copy
  * [x] Delete
  * [x] Download
  * [x] Upload
  * [x] Create folder
* Offline API
  * [x] List tasks
  * [x] Create url task
  * [x] Create task
    * [x] Create URL task
    * [x] Create torrent task
  * [x] Delete tasks
  * [x] Clear tasks
* Other
  * [X] Captcha

# License

![WTFPL](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-2.png)