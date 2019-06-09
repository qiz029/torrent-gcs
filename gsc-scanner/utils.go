package scanner

import "log"

var fileVisited map[string]string

func CheckIfDownloaded(filename string) bool {
	if _, ok := fileVisited[filename]; ok {
		return false
	} else {
		log.Printf("%v has not been downloaded", filename)
		fileVisited[filename] = "visited"
		return true
	}
}
