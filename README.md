# ElevenGo

![Version](https://img.shields.io/badge/release-v0.2.1-brightgreen?style=flat-square)
[![Reference](https://img.shields.io/badge/Go-Reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/deadblue/elevengo)
![License](https://img.shields.io/:License-MIT-green.svg?style=flat-square)

An API client of 115 Cloud Storage Service.

v0.2.x is in process.

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

  files := make([]*elevengo.File, 10)
  for cursor := new(elevengo.FileCursor); cursor.HasMore(); {
    n, err := agent.FileList("0", cursor, files)
    if err != nil {
      log.Fatalf("List file failed: %s", err.Error())
    }
    for i := 0; i < n; i++ {
      log.Printf("File: %#v", files[i])
    }
  }  
}
```

More examples can be found in [reference](https://pkg.go.dev/github.com/deadblue/elevengo).

## Features

* Login
  * [x] Import credential from cookies
  * [x] Login via QRCode
  * [x] Get signed-in user information
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
  * [x] Get Video data
  * [X] Get Image URL
* Offline
  * [x] List tasks
  * [x] Create task by URL
  * [x] Delete tasks
  * [x] Clear tasks
* Other
  * [x] Captcha

## License

MIT
