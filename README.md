# ElevenGo

![Version](https://img.shields.io/badge/release-v0.7.2-brightgreen?style=flat-square)
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

  // Iterate files under specific directory
  if it, err := agent.FileIterate("dirId"); err != nil {
    log.Fatalf("Iterate files error: %s", err)
  } else {
    log.Printf("Total files: %d", it.Count())
    for index, file := range it.Items() {
      log.Printf("%d => %#v", index, file)
    }
  }
  
}
```

</details>

<details>

<summary>Low-level API</summary>

```go
package main

import (
    "context"
    "log"

    "github.com/deadblue/elevengo"
    "github.com/deadblue/elevengo/lowlevel/api"
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
  if err = llc.CallApi(spec, context.Background()); err != nil {
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
