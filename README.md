# ElevenGo ![](https://img.shields.io/badge/status-WIP-green.svg)

An API wrapper for 115 NetDisk Service in Golang.

# Exmaple

Since the API has a lot of changes, example code may be changed in future.

```
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

// TODO: Call some API
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