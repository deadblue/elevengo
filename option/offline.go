package option

type OfflineOption interface {
	isOfflineOption()
}

type OfflineDeleteOption bool

func (o OfflineDeleteOption) isOfflineOption() {}

const OfflineDeleteFiles = OfflineDeleteOption(true)

type OfflineSaveOption string

func (o OfflineSaveOption) isOfflineOption() {}

func OfflineSaveAt(dirId string) OfflineOption {
	return OfflineSaveOption(dirId)
}
