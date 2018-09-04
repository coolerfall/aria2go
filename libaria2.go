// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package aria2go

/*
 #cgo CXXFLAGS: -std=c++11 -I./aria2-lib/include
 #cgo LDFLAGS: -L./aria2-lib/lib
 #cgo LDFLAGS: -laria2 -lssh2 -lcrypto -lssl -lcares -lsqlite3 -lz -lexpat
 #include "aria2_c.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// Type definition for lib aria2, it holds a notifier.
type Aria2 struct {
	notifier Notifier
}

// NewAria2 creates a new instance of aria2.
func NewAria2() *Aria2 {
	return NewAria2WithOptions(nil)
}

// NewAria2WithOptions creates a new instance of aira2 with global options.
// See `ChangeGlobalOptions` also.
func NewAria2WithOptions(options Options) *Aria2 {
	a := &Aria2{
		notifier: newDefaultNotifier(),
	}
	C.init(C.ulong(uintptr(unsafe.Pointer(a))), C.CString(a.fromOptions(options)))
	return a
}

// Start the aria2 to keep running. Note this will block current thread.
func (a *Aria2) Start() {
	C.start()
}

// SetNotifier sets notifier to receive download notification from aria2.
func (a *Aria2) SetNotifier(notifier Notifier) {
	if notifier == nil {
		return
	}
	a.notifier = notifier
}

// AddUri adds a new download. The uris is an array of HTTP/FTP/SFTP/MetaInfo
// URIs (strings) pointing to the same resource. When adding MetaInfo Magnet
// URIs, uris must have only one element and it should be MetaInfo Magnet URI.
func (a *Aria2) AddUri(uri string) (gid string, err error) {
	ret := C.addUri(C.CString(uri))
	if ret == 0 {
		return "", errors.New("add uri failed")
	}
	return fmt.Sprintf("%x", uint64(ret)), nil
}

// ParseTorrent parses torrent file into torrent information. Aria2 will not
// download. This will return info hash, all files in torrent.
func (a *Aria2) ParseTorrent(filepath string) (*BitTorrentInfo, error) {
	ret := C.parseTorrent(C.CString(filepath))
	if ret == nil {
		return nil, errors.New("no data in torrent file")
	}

	// convert info hash to hex string
	infoHash := fmt.Sprintf("%x", []byte(C.GoString(ret.infoHash)))

	// retrieve BitTorrent meta information
	var metaInfo = MetaInfo{}
	mi := ret.metaInfo
	if mi != nil {
		announceList := strings.Split(C.GoString(mi.announceList), ";")
		metaInfo = MetaInfo{
			Name:         C.GoString(mi.name),
			Comment:      C.GoString(mi.comment),
			CreationUnix: int64(mi.creationUnix),
			AnnounceList: announceList,
		}
	}

	return &BitTorrentInfo{
		InfoHash: infoHash,
		MetaInfo: metaInfo,
		Files:    a.parseFiles(ret.files, ret.numFiles),
	}, nil
}

// AddTorrent adds a MetaInfo download with given torrent file path.
// This will return gid and files in torrent file if add successfully.
// User can choose specified files to download, change directory and so on.
func (a *Aria2) AddTorrent(filepath string, options Options) (gid string, err error) {
	ret := C.addTorrent(C.CString(filepath), C.CString(a.fromOptions(options)))
	if ret == 0 {
		return "", errors.New("add torrent failed")
	}
	return fmt.Sprintf("%x", uint64(ret)), nil
}

// ChangeOptions can change the options for aria2. See available options in
// https://aria2.github.io/manual/en/html/aria2c.html#input-file.
func (a *Aria2) ChangeOptions(gid string, options Options) error {
	if !C.changeOptions(a.hexToGid(gid), C.CString(a.fromOptions(options))) {
		return errors.New("change options error")
	}

	return nil
}

// GetOptions gets all options for given gid.
func (a *Aria2) GetOptions(gid string) Options {
	cOptions := C.getOptions(a.hexToGid(gid))
	if cOptions == nil {
		return make(Options)
	}

	return a.toOptions(C.GoString(cOptions))
}

// ChangeGlobalOptions changes global options. See available options in
// https://aria2.github.io/manual/en/html/aria2c.html#input-file except for
// `checksum`, `index-out`, `out`, `pause` and `select-file`.
func (a *Aria2) ChangeGlobalOptions(options Options) error {
	if !C.changeGlobalOptions(C.CString(a.fromOptions(options))) {
		return errors.New("change global options error")
	}

	return nil
}

// GetGlobalOptions gets all global options of aria2.
func (a *Aria2) GetGlobalOptions() Options {
	return a.toOptions(C.GoString(C.getGlobalOptions()))
}

// Pause pauses an active download for given gid. The status of the download
// will become `DOWNLOAD_PAUSED`. Use `Resume` to restart download.
func (a *Aria2) Pause(gid string) bool {
	return bool(C.pause(a.hexToGid(gid)))
}

// Resume resumes an paused download for given gid.
func (a *Aria2) Resume(gid string) bool {
	return bool(C.resume(a.hexToGid(gid)))
}

// Remove removes download no matter what status it was. This will stop
// downloading and stop seeding(for torrent).
func (a *Aria2) Remove(gid string) bool {
	return bool(C.removeDownload(a.hexToGid(gid)))
}

// GetDownloadInfo gets current download information for given gid.
func (a *Aria2) GetDownloadInfo(gid string) DownloadInfo {
	ret := C.getDownloadInfo(a.hexToGid(gid))
	if ret == nil {
		return DownloadInfo{}
	}

	return DownloadInfo{
		Status:         int(ret.status),
		TotalLength:    int64(ret.totalLength),
		BytesCompleted: int64(ret.bytesCompleted),
		BytesUpload:    int64(ret.uploadLength),
		DownloadSpeed:  int(ret.downloadSpeed),
		UploadSpeed:    int(ret.uploadSpeed),
		Files:          a.parseFiles(ret.files, ret.numFiles),
	}
}

// fromOptions converts `Options` to string with ';' separator.
func (a *Aria2) fromOptions(options Options) string {
	if options == nil {
		return ""
	}

	var cOptions string
	for k, v := range options {
		cOptions += k + ";"
		cOptions += v + ";"
	}

	return strings.TrimSuffix(cOptions, ";")
}

// fromOptions converts options string with ';' separator to `Options`.
func (a *Aria2) toOptions(cOptions string) Options {
	coptions := strings.Split(strings.TrimSuffix(cOptions, ";"), ";")
	var options = make(Options)
	var index int
	for index = 0; index < len(coptions); index += 2 {
		options[coptions[index]] = coptions[index+1]
	}

	return options
}

// hexToGid convert hex to uint64 type gid.
func (a *Aria2) hexToGid(hex string) C.ulong {
	id, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return 0
	}
	return C.ulong(id)
}

// parseFiles parses all files information from aria2.
func (a *Aria2) parseFiles(filesPointer *C.struct_FileInfo, length C.int) (files []File) {
	cfiles := (*[1 << 30]C.struct_FileInfo)(unsafe.Pointer(filesPointer))[:length:length]
	if cfiles == nil {
		return
	}

	for _, f := range cfiles {
		files = append(files, File{
			Index:           int(f.index),
			Length:          int64(f.length),
			CompletedLength: int64(f.completedLength),
			Name:            C.GoString(f.name),
			Selected:        bool(f.selected),
		})
	}

	return
}

//export notifyEvent
//noinspection GoUnusedFunction
func notifyEvent(ariagoPointer uint64, id uint64, event int) {
	a := (*Aria2)(unsafe.Pointer(uintptr(ariagoPointer)))
	if a == nil || a.notifier == nil {
		return
	}

	// convert id to hex string
	gid := fmt.Sprintf("%x", uint64(id))

	switch event {
	case onStart:
		a.notifier.OnStart(gid)
	case onPause:
		a.notifier.OnPause(gid)
	case onStop:
		a.notifier.OnStop(gid)
	case onComplete:
		a.notifier.OnComplete(gid)
	case onError:
		a.notifier.OnError(gid)
	case onBTComplete:
		a.notifier.OnComplete(gid)
	}
}
