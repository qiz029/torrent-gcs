package main

import (
	"log"
	"net"
	"time"

	"github.com/anacrolix/torrent"
	scanner "github.com/qiz029/torrent-gcs/gsc-scanner"
)

const (
	videoFileLoc    = "~/torrents-store/videos"
	torrentsFileLoc = "~/torrents-store/torrents"
)

func main() {
	log.Printf("start to run torrent-gcs")
	clientConf := torrent.NewDefaultClientConfig()
	clientConf.Debug = false
	clientConf.Seed = false // seed after downloaded
	clientConf.DataDir = videoFileLoc
	var publicIP net.IP
	clientConf.PublicIp4 = publicIP
	clientConf.PublicIp6 = publicIP

	client, err := torrent.NewClient(clientConf)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	defer client.Close()

	// From here we need 2 threads
	// 1. restApi to take magnet download
	// 2. scan the GSC to download
	forever := make(chan bool)
	go scanner.GCSScanner(client, torrentsFileLoc, 30*time.Second)

	<-forever
}
