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
	Announces       [][]string `json:"announceList"`
	Comment         string     `json:"comment"`
	CreationSeconds int        `json:"creationDate"`
	Mode            string     `json:"mode"`
	VerifiedLength  int64      `json:"verifiedLength,string"`
	VerifyPending   bool       `json:"verifyIntegrityPending,string"`
}

// Type definition for bit torrent detail information.
type BitTorrentInfo struct {
	Hash     string     `json:"infoHash"`
	Seeders  int        `json:"numSeeders,string"`
	IsSeeder bool       `json:"seeder,string"`
	Torrent  BitTorrent `json:"bittorrent"`
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
	PeerId        string `json:"peerId"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	DownloadSpeed string `json:"downloadSpeed"`
	UploadSpeed   string `json:"uploadSpeed"`
	IsSeeder      bool   `json:"seeder,string"`
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
