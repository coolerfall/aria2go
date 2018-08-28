// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package aria2go

// Type definition for download information.
type DownloadInfo struct {
	Status         int
	TotalLength    int64
	BytesCompleted int64
	BytesUpload    int64
	DownloadSpeed  int
	UploadSpeed    int
	BitField       string
}

// Type definition for BitTorrent meta information.
type MetaInfo struct {
	Name         string
	AnnounceList []string
	Comment      string
	CreationUnix int64
	Mode         string
}

// Type definition for file in torrent.
type File struct {
	Index    int
	Name     string
	Length   int64
	Selected bool
}

// Type definition for BitTorrent detail information.
type BitTorrentInfo struct {
	InfoHash string
	MetaInfo MetaInfo
	Files    []File
	Seeders  int
	IsSeeder bool
}

type Options map[string]string

// Type definition for download event, this will keep the same with aria2.
const (
	onStart = iota + 1
	onPause
	onStop
	onComplete
	onError
	onBTComplete
)
