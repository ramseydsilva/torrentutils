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
	Length int
	Md5sum string
	Path   []string
}

type Info struct {
	PieceLength int "piece length"
	Pieces      string
	Private     int
	Name        string
	Length      int
	Md5sum      string
	Files       []File
}

type MetaInfo struct {
	Info         Info
	Announce     string
	AnnounceList []string "announce-list"
	CreationDate int64    "creation date"
	Comment      string
	CreatedBy    string "created by"
	Encoding     string
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
	const pieceLength = 2000000 // we setle for ~ 800KB

	fileInfo, err := file.Stat()
	checkError(err)

	info := &Info{
		PieceLength: pieceLength,
		Private:     0,
	}

	if oneFile {
		if clf.Name == "" {
			info.Name = fileInfo.Name()
		} else {
			info.Name = clf.Name
		}
		info.Length = int(fileInfo.Size())
		info.Md5sum, info.Pieces = getMd5SumAndPieces(file, fileInfo.Size(), pieceLength)
	} else {
		log.Fatalf("Don't support multiple files, yet")
	}

	metaInfo := &MetaInfo{
		Info:         *info,
		Announce:     clf.Announce,
		AnnounceList: strings.Split(clf.AnnounceList, ","),
		CreationDate: time.Now().Unix(),
		Comment:      clf.Comment,
		CreatedBy:    clf.CreatedBy,
		Encoding:     clf.Encoding,
	}

	torrentFile, err := os.Create(strings.Replace(filename, fileInfo.Name(), info.Name, -1) + ".torrent")
	defer torrentFile.Close()
	checkError(err)
	bencode.Marshal(torrentFile, *metaInfo)

	return torrentFile
}
