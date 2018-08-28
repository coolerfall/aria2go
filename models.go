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

// Type definition for bit torrent.
type BitTorrent struct {
	Announces       [][]string
	Comment         string
	CreationSeconds int
	Mode            string
	VerifiedLength  int64
	VerifyPending   bool
}

// Type definition for bit torrent detail information.
type BitTorrentInfo struct {
	Hash     string
	Seeders  int
	IsSeeder bool
	Torrent  BitTorrent
}

// Type definition for file in torrent.
type File struct {
	Index    int
	Name     string
	Length   int64
	Selected bool
}

type Options map[string]string

// Type definition for peer of bit torrent.
type Peer struct {
	PeerId        string
	Ip            string
	Port          string
	DownloadSpeed string
	UploadSpeed   string
	IsSeeder      bool
}

// Type definition for download event, this will keep the same with aria2.
const (
	onStart = iota + 1
	onPause
	onStop
	onComplete
	onError
	onBTComplete
)
