# ElevenGo

[![GoDoc](https://godoc.org/github.com/deadblue/elevengo?status.svg)](https://godoc.org/github.com/deadblue/elevengo)

A Golang API wrapper for 115 NetDisk Service.

## Example

> You can found more example codes godoc: [https://godoc.org/github.com/deadblue/elevengo](https://godoc.org/github.com/deadblue/elevengo).

```go
import "github.com/deadblue/elevengo"

// Create agent
agent = elevengo.Default()

// Import credentials to login
cr = &elevengo.Credentials{
    UID: "",
    CID: "",
    SEID: "",
}
if err := agent.ImportCredentials(cr); err != nil {
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

## Features

* Login
  * [x] Import credentials from cookies
  * [x] Login via QRCode
  * [ ] ~~Login via Account/Password~~ (No idea)
* File API
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
* Offline API
  * [x] List tasks
  * [x] Create URL task(s)
  * [x] Delete tasks
  * [x] Clear tasks
* Other
  * [x] Captcha

## TODO list

* Current version:
  * Print some logs via Logger interface.
* Next version:
  * Handle more upstream errors.
  * Implement download/upload method, with progress echo.

## License

MIT