// Copyright (c) 2018-present Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"log"
	"time"

	"anbillon.com/aria2go"
)

func main() {
	a := aria2go.NewAria2WithOptions(aria2go.Options{
		"dir": "test/",
	})
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

	for k, v := range a.GetOptions(gid) {
		log.Printf("%v: %v", k, v)
	}

	for {
		time.Sleep(time.Second)
		di := a.GetDownloadInfo(gid)
		for _, f := range di.Files {
			log.Printf("index: %v, name: %v, selected: %v, completed: %v, len: %v",
				f.Index, f.Name, f.Selected, f.CompletedLength, f.Length)
		}
		log.Printf("download speed: %vKib/s", di.DownloadSpeed/1024.0)
		log.Printf("bytes completed: %vM", float64(di.BytesCompleted)/1024.0/1024.0)
	}
}
