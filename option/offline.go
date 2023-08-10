package option

type OfflineDeleteOption interface {
	isOfflineDeleteOption()
}

type OfflineDeleteFilesOfTasks bool

func (o OfflineDeleteFilesOfTasks) isOfflineDeleteOption() {}

type OfflineAddOption interface {
	isOfflineAddOption()
}

type OfflineSaveDownloadedFileTo string

func (o OfflineSaveDownloadedFileTo) isOfflineAddOption() {}
