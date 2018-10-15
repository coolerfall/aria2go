// Copyright (c) 2018 Anbillon Team (anbillonteam@gmail.com).

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"anbillon.com/aria2go"
	"gopkg.in/cheggaaa/pb.v1"
)

func main() {
	flag.Parse()
	flag.Usage = func() {
		log.Printf("Usage: aria2goc uri/torrent")
	}
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	a := aria2go.NewAria2(aria2go.Config{
		Options: aria2go.Options{
			"dir":          "./data/",
			"save-session": "./data/aria2go.session",
		},
	})

	input := flag.Arg(0)
	isMagnet := strings.HasPrefix(input, "magnet")
	if strings.HasPrefix(input, "http") || isMagnet {
		_, err := a.AddUri(input, nil)
		if err != nil {
			flag.Usage()
			return
		}
	} else if strings.HasSuffix(input, "torrent") {
		_, err := a.AddTorrent(input, nil)
		if err != nil {
			flag.Usage()
			return
		}
	} else {
		log.Printf("not supported uri or file")
		return
	}

	go func() {
		a.Run()
	}()

	startNotif := make(chan string)
	completeNotif := make(chan bool)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGHUP)

		select {
		case <-quit:
			completeNotif <- true
		}
	}()
	a.SetNotifier(newAria2goNotifier(startNotif, completeNotif))

	var gid string
	var bar *pb.ProgressBar
	var count int
	ticker := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case id := <-startNotif:
			{
				gid = id
				bar = createProgressBar(a, gid)
				if bar == nil {
					log.Printf("fetch download information error")
					completeNotif <- true
					return
				} else {
					bar.Start()
				}
			}
		case <-ticker.C:
			showProgress(a, gid, bar)
		case shutdown := <-completeNotif:
			bar.Finish()
			if !shutdown && isMagnet && count < 1 {
				continue
			}
			count ++
			ticker.Stop()
			os.Exit(1)
		}
	}
}

func createProgressBar(a *aria2go.Aria2, gid string) (bar *pb.ProgressBar) {
	var retry int
	for {
		if retry >= 60 {
			return nil
		}
		di := a.GetDownloadInfo(gid)
		if di.TotalLength <= 0 {
			retry ++
			time.Sleep(time.Second)
			continue
		}

		bar = pb.New64(di.TotalLength)
		bar.SetUnits(pb.U_BYTES)
		bar.ShowElapsedTime = false
		bar.ShowTimeLeft = true
		bar.ShowPercent = true
		bar.ShowSpeed = true
		name := di.MetaInfo.Name
		if len(di.MetaInfo.Name) != 0 {
			name = di.MetaInfo.Name
		} else if len(di.Files) > 0 {
			name = di.Files[0].Name
		}
		bar.Prefix(name)

		return
	}
}

func showProgress(a *aria2go.Aria2, gid string, pb *pb.ProgressBar) {
	if pb == nil {
		return
	}
	di := a.GetDownloadInfo(gid)
	pb.Set64(di.BytesCompleted)
}

type Aria2gocNotifier struct {
	start    chan string
	complete chan bool
}

func newAria2goNotifier(start chan string, complete chan bool) aria2go.Notifier {
	return Aria2gocNotifier{
		start:    start,
		complete: complete,
	}
}

func (n Aria2gocNotifier) OnStart(gid string) {
	log.Printf("on start %v", gid)
	n.start <- gid
}

func (n Aria2gocNotifier) OnPause(gid string) {
	log.Printf("on pause: %v", gid)
}

func (n Aria2gocNotifier) OnStop(gid string) {
	log.Printf("on stop: %v", gid)
}

func (n Aria2gocNotifier) OnComplete(gid string) {
	log.Printf("on complete: %v", gid)
	n.complete <- false
}

func (n Aria2gocNotifier) OnError(gid string) {
	log.Printf("on error: %v", gid)
}
