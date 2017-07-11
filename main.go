package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var (
	filepath = flag.String("filepath", "example.mkv", "Filepath whose subtitles will be downloaded")
	language = flag.String("language", "en", "Provide language of subtitle")
)

func main() {
	flag.Parse()
	hash := getHash(*filepath)
	userAgent := getUA()
	downloadSubtitles(hash, userAgent, *language, *filepath)
}

func createSubtitleFilepath(filepath string) string {
	fileLength := len(filepath)
	fileExt := path.Ext(filepath)
	return filepath[:fileLength-len(fileExt)] + ".srt"
}

func downloadSubtitles(hash, userAgent, language, filepath string) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://api.thesubdb.com/?action=download&hash=%s&language=%s", hash, language), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	subtitleFilepath := createSubtitleFilepath(filepath)
	subtitleFile, err := os.OpenFile(subtitleFilepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer subtitleFile.Close()
	_, err = io.Copy(subtitleFile, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func getUA() string {
	clientName := "Surabhi"
	clientVersion := "1.0"
	clientURL := "https://github.com/surabhigupta412"
	return fmt.Sprintf("SubDB/1.0 (%s/%s; %s)", clientName, clientVersion, clientURL)
}

func getHash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var readSize int64 = 64 * 1024
	hash := md5.New()
	_, err = io.CopyN(hash, file, readSize)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Seek(-readSize, os.SEEK_END)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.CopyN(hash, file, readSize)
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
