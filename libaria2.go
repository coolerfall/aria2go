// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package aria2go

/*
 #cgo CXXFLAGS: -std=c++11
 #cgo LDFLAGS: -lcrypto -lgcrypt
 #cgo LDFLAGS: -L ./aria2-lib/lib -laria2 -lssh2 -lcares -lsqlite3 -lz -lexpat -lssl
 #include "aria2_c.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

// Type definition for lib aria2, it holds a notifier.
type Aria2 struct {
	notifier Notifier
}

// NewAria2 creates a new instance of aria2.
func NewAria2() *Aria2 {
	a := &Aria2{
		notifier: NewDefaultNotifier(),
	}
	C.init(C.ulong(uintptr(unsafe.Pointer(a))))
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

// AddUri adds a new download. The uris is an array of HTTP/FTP/SFTP/BitTorrent
// URIs (strings) pointing to the same resource. When adding BitTorrent Magnet
// URIs, uris must have only one element and it should be BitTorrent Magnet URI.
func (a *Aria2) AddUri(uri string) (gid string, err error) {
	ret := C.addUri(C.CString(uri))
	if ret == 0 {
		return "", errors.New("add uri failed")
	}
	return fmt.Sprintf("%x", uint64(ret)), nil
}

// ParseTorrent parses torrent file into torrent information.
// This will return all files and their size.
func (a *Aria2) ParseTorrent(filepath string) (files []File, err error) {
	ret := C.parseTorrent(C.CString(filepath))
	if ret == nil {
		return nil, errors.New("no data in torrent file")
	}

	length := ret.totalFile
	cfiles := (*[1 << 30]C.struct_FileInfo)(unsafe.Pointer(ret.files))[:length:length]
	for _, f := range cfiles {
		files = append(files, File{
			Index:    int(f.index),
			Length:   int64(f.length),
			Name:     C.GoString(f.name),
			Selected: bool(f.selected),
		})
	}

	return
}

// AddTorrent adds a BitTorrent download with given torrent file path.
// This will return gid and files in torrent file if add successfully.
// User can choose specified files to download, change directory and so on.
func (a *Aria2) AddTorrent(filepath string, options Options) (gid string, err error) {
	var cOptions string
	for k, v := range options {
		cOptions += k
		cOptions += ","
		cOptions += v
	}

	ret := C.addTorrent(C.CString(filepath), C.CString(cOptions))
	if ret == 0 {
		return "", errors.New("add torrent failed")
	}
	return fmt.Sprintf("%x", uint64(ret)), nil
}

// ChangeOptions can change the options for aria2. See available options in
// https://aria2.github.io/manual/en/html/aria2c.html#input-file.
func (a *Aria2) ChangeOptions(gid string, options Options) error {
	var cOptions string
	for k, v := range options {
		cOptions += k
		cOptions += ","
		cOptions += v
	}

	if !C.changeOptions(a.hexToGid(gid), C.CString(cOptions)) {
		return errors.New("change option error")
	}

	return nil
}

// Pause pauses an active downloading for given gid. The status of the download
// will become `DOWNLOAD_PAUSED`. Use `Resume` to restart download.
func (a *Aria2) Pause(gid string) bool {
	return bool(C.pause(a.hexToGid(gid)))
}

// Resume resumes an paused downloading for given gid.
func (a *Aria2) Resume(gid string) bool {
	return bool(C.resume(a.hexToGid(gid)))
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
	}
}

// hexToGid convert hex to uint64 type gid.
func (a *Aria2) hexToGid(hex string) C.ulong {
	id, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return 0
	}
	return C.ulong(id)
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
