# ElevenGo ![](https://img.shields.io/badge/status-WIP-green.svg)

An API wrapper for 115 NetDisk Service in Golang.

# Usage

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

// Call APIs
// TODO
```

# Features

* Login
 * [x] Import credentials from cookies
 * [ ] Login via QRCode
 * [ ] Login via Account/Password (No idea now)
* File API
 * [x] List/Search
 * [x] Create folder
 * [x] Rename/Move/Copy
 * [x] Delete
 * [x] Download
* Offline API
 * [x] List tasks
 * [x] Create task
 * [x] Delete/Clear tasks
* Other
 * [X] Captcha

# License

![WTFPL](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-2.png)