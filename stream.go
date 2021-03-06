package main

import (
	"log"
	"time"

	"github.com/webver/vdk/format/rtsp"
)

func serveStreams() {
	for k, v := range Config.Streams {
		go func(name, url string) {
			stopHlsCast := make(chan bool, 1)
			for {
				log.Println(name, "connect", url)
				rtsp.DebugRtsp = true
				rtsp.DebugRtp = false
				session, err := rtsp.Dial(url)
				if err != nil {
					log.Println(name, err)
					time.Sleep(5 * time.Second)
					continue
				}
				session.RtpKeepAliveTimeout = time.Duration(10 * time.Second)
				if err != nil {
					log.Println(name, err)
					time.Sleep(5 * time.Second)
					continue
				}
				codec, err := session.Streams()
				if err != nil {
					log.Println(name, err)
					time.Sleep(5 * time.Second)
					continue
				}
				Config.codecAdd(name, codec)
				Config.updateStatus(name, true)
				Config.startHlsCast(name, stopHlsCast)
				for {
					pkt, err := session.ReadPacket()
					if err != nil {
						log.Println(name, err)
						break
					}
					Config.cast(name, pkt)
				}
				session.Close()
				stopHlsCast <- true
				Config.updateStatus(name, false)
				log.Println(name, "reconnect wait 5s")
				time.Sleep(5 * time.Second)
			}
		}(k, v.URL)
	}
}
