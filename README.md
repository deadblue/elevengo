# ElevenGo ![](https://img.shields.io/badge/status-WIP-green.svg)

An API wrapper for 115 NetDisk Service in Golang.

# Example

Since the API has a lot of changes, example code may be changed in future.

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
        // TODO: handle the files
    }
}

// List offline tasks.
for cursor := elevengo.OfflineCursor(); cursor.HasMore(); cursor.Next() {
    if tasks, err := agent.OfflineList(cursor); err != nil {
        panic(err)
    } else {
        // TODO: handle the tasks
    }
}

```

# Features

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
  * [ ] Stat (W.I.P.)
  * [x] Download
  * [x] Upload
* Offline API
  * [x] List tasks
  * [x] Create URL task(s)
  * [x] Delete tasks
  * [x] Clear tasks
* Other
  * [X] Captcha

# Following works to do:

* ~~Play games.(NO!)~~
* Handle more upstream errors.
* Re-design error system, merge all errors into one type.
* Implement download/upload method, with progress echo.
* Print logging throuhg Logger interface.
* Add more docs and example codes.

# License

MIT