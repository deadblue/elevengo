package option

type FileListOptions struct {
	/*Predefined file type:
	- 0: All
	- 1: Document
	- 2: Image
	- 3: Audio
	- 4: Video
	- 5: Archive
	- 6: Sofrwore
	*/
	Type int
	// File extension with leading-dot, e.g.: mkv
	ExtName string
}

func (o *FileListOptions) ShowAll() *FileListOptions {
	o.Type = 0
	o.ExtName = ""
	return o
}

func (o *FileListOptions) OnlyDocument() *FileListOptions {
	o.Type = 1
	o.ExtName = ""
	return o
}

func (o *FileListOptions) OnlyImage() *FileListOptions {
	o.Type = 2
	o.ExtName = ""
	return o
}

func (o *FileListOptions) OnlyAudio() *FileListOptions {
	o.Type = 3
	o.ExtName = ""
	return o
}

func (o *FileListOptions) OnlyVideo() *FileListOptions {
	o.Type = 4
	o.ExtName = ""
	return o
}

func (o *FileListOptions) OnlyArchive() *FileListOptions {
	o.Type = 5
	o.ExtName = ""
	return o
}

func (o *FileListOptions) OnlySoftware() *FileListOptions {
	o.Type = 6
	o.ExtName = ""
	return o
}

func (o *FileListOptions) OnlyExtension(extName string) *FileListOptions {
	o.Type = -1
	o.ExtName = extName
	return o
}

func FileList() *FileListOptions {
	return (&FileListOptions{}).ShowAll()
}
