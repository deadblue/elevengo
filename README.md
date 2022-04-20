# ElevenGo

![Version](https://img.shields.io/badge/release-v0.2.1-brightgreen?style=flat-square)
[![Reference](https://img.shields.io/badge/Go-Reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/deadblue/elevengo)
![License](https://img.shields.io/:License-MIT-green.svg?style=flat-square)

A Golang API wrapper for 115 Cloud Storage Service.

## Example

```go
import "github.com/deadblue/elevengo"

// Create agent
agent = elevengo.Default()

// Import credentials to login
cr := &elevengo.Credential{
    UID: "",
    CID: "",
    SEID: "",
}
if err := agent.CredentialImport(cr); err != nil {
    panic(err)
}

// List files under root.
for cursor := elevengo.FileCursor(); cursor.HasMore(); cursor.Next() {
    if files, err := agent.FileList("0", cursor); err != nil {
        panic(err)
    } else {
        // TODO: deal with the files
    }
}
```

You can find more examples from [reference](https://pkg.go.dev/github.com/deadblue/elevengo).

## Features

* Login
  * [x] Import credential from cookies
  * [x] Login via QRCode
  * [x] Get signed in user information
* File
  * [x] List
  * [x] Search
  * [x] Rename
  * [x] Move
  * [x] Copy
  * [x] Delete
  * [x] Mkdir
  * [x] Stat
  * [x] Storage Stat
  * [x] Download
  * [x] Upload
  * [x] Video HLS
  * [X] Image URL
* Offline
  * [x] List tasks
  * [x] Create URL task(s)
  * [x] Delete tasks
  * [x] Clear tasks
* Other
  * [x] Captcha

## TODO list

* Handle more upstream errors.
* Caller can switch upstream API between HTTP/HTTPS.

## License

MIT
