package option

type FileOption interface {
	isFileOption()
}

type FileKindOption int

func (o FileKindOption) isFileOption() {}

type FileStarOption bool

func (o FileStarOption) isFileOption() {}

type FileLabelOption string

func (o FileLabelOption) isFileOption() {}

func FileSearchByLabel(labelId string) FileLabelOption {
	return FileLabelOption(labelId)
}

type FileKeyworkOption string

func (o FileKeyworkOption) isFileOption() {}

func FileSearchByKeyword(keyword string) FileKeyworkOption {
	return FileKeyworkOption(keyword)
}

const (
	FileKindAll FileKindOption = iota
	FileKindDocument
	FileKindImage
	FileKindAudio
	FileKindVideo
	FileKindArchive
	FileKindApplication

	FileStared FileStarOption = true
)
