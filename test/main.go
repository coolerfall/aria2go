// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"anbillon.com/aria2go"
	"log"
	"time"
)

func main() {
	a := aria2go.NewAria2()
	go func() {
		a.Start()
	}()
	// gid, _ := a.AddUri("http://mirrors.evowise.com/archlinux/iso/2018.08.01/archlinux-2018.08.01-x86_64.iso")
	// log.Printf("gid: %v", gid)
	btInfo, err := a.ParseTorrent("test/test.torrent")
	if err != nil {
		return
	}
	log.Printf("info hash: %v", btInfo.InfoHash)
	log.Printf("name: %v", btInfo.MetaInfo.Name)
	log.Printf("announce list: %v", btInfo.MetaInfo.AnnounceList)
	for _, f := range btInfo.Files {
		log.Printf("file %v", f)
	}

	gid, err := a.AddTorrent("test/test.torrent", nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	for {
		time.Sleep(time.Second)
		di := a.GetDownloadInfo(gid)
		log.Printf("download speed: %vKib/s", di.DownloadSpeed/1024.0)
		log.Printf("bytes completed: %vM", float64(di.BytesCompleted)/1024.0/1024.0)
	}
}
