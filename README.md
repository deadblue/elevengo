# ElevenGo

[![GoDoc](https://godoc.org/github.com/deadblue/elevengo?status.svg)](https://godoc.org/github.com/deadblue/elevengo)
![](https://img.shields.io/badge/release-v0.1.3-brightgreen?style=flat-square)
![](https://img.shields.io/badge/develop-v0.1.4-orange?style=flat-square)

A Golang API wrapper for 115 Cloud Storage Service.

## Example

```go
import "github.com/deadblue/elevengo"

// Create agent
agent = elevengo.Default()

// Import credentials to login
cr := elevengo.Credential{
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

You can find more example on [GoDoc](https://godoc.org/github.com/deadblue/elevengo).

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
* Offline
  * [x] List tasks
  * [x] Create URL task(s)
  * [x] Delete tasks
  * [x] Clear tasks
* Other
  * [x] Captcha

## TODO list

* Handle more upstream errors.
* Caller can swtich upstream API between HTTP/HTTPS.
* Implement download/upload method, with progress echo.

## License

MIT