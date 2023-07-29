package option

type FileOption interface {
	isFileOption()
}

type FileKindOption int

func (o FileKindOption) isFileOption() {}

const (
	FileKindAll         FileKindOption = iota
	FileKindDocument
	FileKindImage
	FileKindAudio
	FileKindVideo
	FileKindArchive
	FileKindApplication
)
