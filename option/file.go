package option

type FileFilterOption interface {
	isFileFilterOption()
}

type FileFilterByKind int

const (
	FileKindAll FileFilterByKind = iota
	FileKindDocument
	FileKindImage
	FileKindAudio
	FileKindVideo
	FileKindArchive
	FileKindApplication
)
