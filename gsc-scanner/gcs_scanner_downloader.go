package scanner

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/dustin/go-humanize"
	"github.com/gosuri/uiprogress"
)

var progress = uiprogress.New()

func GCSScanner(client *torrent.Client, torrentLoc string, timeInterval time.Duration) {
	log.Printf("start to scan GCS bucket %v every %v", torrentLoc, timeInterval)
	files := make([]string, 0)

	for {
		err := filepath.Walk(torrentLoc, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("met error %v", err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".torrent" {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			log.Printf("GSC scan error met %v", err)
			continue
		}

		go DownloadFile(client, files)
		time.Sleep(timeInterval)
	}
}

func DownloadFile(client *torrent.Client, torrentFiles []string) {
	var t *torrent.Torrent
	t = nil
	for _, v := range torrentFiles {
		metaInfo, err := metainfo.LoadFromFile(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading torrent file %v: %s\n", v, err)
			log.Printf("cannot download because of %v", err)
			return
		}
		t, err = client.AddTorrent(metaInfo)
		if err != nil {
			log.Printf("torrent creation failed because of %v", err)
			return
		}
	}
	// torrentBar(t)
	if t != nil {
		go func() {
			<-t.GotInfo()
			t.DownloadAll()
		}()
	}
}

func torrentBar(t *torrent.Torrent) {
	bar := progress.AddBar(1)
	bar.AppendCompleted()
	bar.AppendFunc(func(*uiprogress.Bar) (ret string) {
		select {
		case <-t.GotInfo():
		default:
			return "getting info"
		}
		if t.Seeding() {
			return "seeding"
		} else if t.BytesCompleted() == t.Info().TotalLength() {
			return "completed"
		} else {
			return fmt.Sprintf("downloading (%s/%s)", humanize.Bytes(uint64(t.BytesCompleted())), humanize.Bytes(uint64(t.Info().TotalLength())))
		}
	})
	bar.PrependFunc(func(*uiprogress.Bar) string {
		return t.Name()
	})
	go func() {
		<-t.GotInfo()
		tl := int(t.Info().TotalLength())
		if tl == 0 {
			bar.Set(1)
			return
		}
		bar.Total = tl
		for {
			bc := t.BytesCompleted()
			bar.Set(int(bc))
			time.Sleep(time.Second)
		}
	}()
}
