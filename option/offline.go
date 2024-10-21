package option

type OfflineAddOptions struct {
	SaveDirId string
}

func (o *OfflineAddOptions) WithSaveDirId(dirId string) *OfflineAddOptions {
	o.SaveDirId = dirId
	return o
}

func OfflineAdd() *OfflineAddOptions {
	return &OfflineAddOptions{}
}

type OfflineDeleteOptions struct {
	DeleteFiles bool
}

func (o *OfflineDeleteOptions) DeleteDownloadedFiles() *OfflineDeleteOptions {
	o.DeleteFiles = true
	return o
}

func OfflineDelete() *OfflineDeleteOptions {
	return &OfflineDeleteOptions{}
}
