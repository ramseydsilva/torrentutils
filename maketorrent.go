package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/jackpal/bencode-go"
)

type File struct {
	length int
	md5sum string
	path   []string
}

type Info struct {
	pieceLength int "piece length"
	pieces      string
	private     int
	name        string
	length      int
	md5sum      string
	files       []File
}

type MetaInfo struct {
	info         Info
	announce     string
	announceList []string "announce-list"
	creationDate int64    "creation date"
	comment      string
	createdBy    string "created by"
	encoding     string
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getMd5SumAndPieces(file *os.File, filesize int64, filechunk uint64) (string, string) {
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()
	pieces := ""
	sha := sha1.New()
	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(float64(filechunk), float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)

		file.Read(buf)
		io.WriteString(hash, string(buf)) // append into the hash
		pieces += string(sha.Sum(buf))
	}
	return hex.EncodeToString(hash.Sum(nil)), pieces
}

func MakeTorrentFile(filename string, clf *CLFlags) *os.File {
	file, err := os.Open(filename)
	defer file.Close()
	checkError(err)

	oneFile := true
	const pieceLength = 256000

	fileInfo, err := file.Stat()
	checkError(err)

	info := &Info{
		pieceLength: pieceLength,
		private:     0,
	}

	if oneFile {
		if clf.Name == "" {
			info.name = fileInfo.Name()
		} else {
			info.name = clf.Name
		}
		info.length = int(fileInfo.Size())
		info.md5sum, info.pieces = getMd5SumAndPieces(file, fileInfo.Size(), pieceLength)
	} else {
		log.Fatalf("Don't support multiple files, yet")
	}

	metaInfo := &MetaInfo{
		info:         *info,
		announce:     clf.Announce,
		announceList: strings.Split(clf.AnnounceList, ","),
		creationDate: time.Now().Unix(),
		comment:      clf.Comment,
		createdBy:    clf.CreatedBy,
		encoding:     clf.Encoding,
	}

	torrentFile, err := os.Create(strings.Replace(filename, fileInfo.Name(), info.name, -1) + ".torrent")
	defer torrentFile.Close()
	checkError(err)
	bencode.Marshal(torrentFile, *metaInfo)

	return torrentFile
}
