# ElevenGo

![Version](https://img.shields.io/badge/release-v0.6.1-brightgreen?style=flat-square)
[![Reference](https://img.shields.io/badge/Go-Reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/deadblue/elevengo)
![License](https://img.shields.io/:License-MIT-green.svg?style=flat-square)

An API client of 115 Cloud Storage Service.

## Example

<details>

<summary>High-level API</summary>

```go
package main

import (
    "github.com/deadblue/elevengo"
    "log"
)

func main()  {
  // Initialize agent
  agent := elevengo.Default()
  // Import credential
  credential := &elevengo.Credential{
    UID: "", CID: "", SEID: "",
  }
  if err := agent.CredentialImport(credential); err != nil {
    log.Fatalf("Import credentail error: %s", err)
  }

  // Get file list
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

</details>

<details>

<summary>Low-level API</summary>

```go
package main

import (
    "github.com/deadblue/elevengo"
    "github.com/deadblue/elevengo/lowlevel/api"
    "log"
)

func main()  {
  // Initialize agent
  agent := elevengo.Default()
  // Import credential
  credential := &elevengo.Credential{
    UID: "", CID: "", SEID: "",
  }
  if err := agent.CredentialImport(credential); err != nil {
    log.Fatalf("Import credentail error: %s", err)
  }

  // Get low-level API client
  llc := agent.LowlevelClient()
  // Init FileList API spec
  spec := (&api.FiieListSpec{}).Init("dirId", 0, 32)
  // Call API
  if err = llc.CallApi(spec); err != nil {
    log.Fatalf("Call API error: %s", err)
  }
  // Parse API result
  for index, file := range spec.Result.Files {
    log.Printf("File: %d => %v", index, file)
  }
  
}
```
</details>

## License

MIT
