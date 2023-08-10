package option

type FileListOption interface {
	isFileListOption()
}

type FileListTypeOption int

func (o FileListTypeOption) isFileListOption() {}

const (
	FileTypeAll FileListTypeOption = iota
	FileTypeDocument
	FileTypeImage
	FileTypeAudio
	FileTypeVideo
	FileTypeArchive
	FileTypeApplication
)
