# ElevenGo

![Version](https://img.shields.io/badge/release-v0.5.2-brightgreen?style=flat-square)
[![Reference](https://img.shields.io/badge/Go-Reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/deadblue/elevengo)
![License](https://img.shields.io/:License-MIT-green.svg?style=flat-square)

An API client of 115 Cloud Storage Service.

## Example

```go
package main

import (
    "github.com/deadblue/elevengo"
    "log"
)

func main()  {
  agent := elevengo.Default()
  credential := &elevengo.Credential{
    UID: "", CID: "", SEID: "",
  }
  if err := agent.CredentialImport(credential); err != nil {
    log.Fatalf("Import credentail error: %s", err)
  }

  it, err := agent.FileIterate("dirId")
  for ; err == nil; err = it.Next() {
    file := &elevengo.File{}
    if err = it.Get(file); err == nil {
      log.Printf("File: %d => %#v", it.Index(), file)
    }
  }
  if !elevengo.IsIteratorEnd(err) {
    log.Fatalf("Iterate files error: %s", err)
  }
}
```

More examples can be found in [reference](https://pkg.go.dev/github.com/deadblue/elevengo).

## Features

* Login
  * [x] Import Credential from Cookies
  * [x] Login via QRCode
  * [x] Get Information signed-in User
* File
  * [x] List
  * [x] Search
  * [x] Rename
  * [x] Move
  * [x] Copy
  * [x] Delete
  * [x] Get Information by ID
  * [x] Stat File
  * [x] Download
  * [x] Upload
  * [x] Make Directory
* Media
  * [x] Get Play URL of Video File
  * [X] Get Directly-Access URL of Image File
* Offline
  * [x] List Tasks
  * [x] Create Task by URL
  * [x] Delete Tasks
  * [x] Clear Tasks

## License

MIT
