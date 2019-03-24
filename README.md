# ElevenGo

An API wrapper for 115 NetDisk Service for Golang

# Usage

```
import "github.com/deadblue/elevengo"

// Get client instance with default config
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

## License

![WTFPL](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-2.png)