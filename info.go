package main

import (
	"bytes"
	"errors"
	"os"

	"github.com/jackpal/bencode-go"
)

func TorrentInfo(filename string) *MetaInfo {
	file, err := os.Open(filename)
	defer file.Close()
	checkError(err)

	m, err := bencode.Decode(file)
	checkError(err)

	metaMap, ok := m.(map[string]interface{})
	if !ok {
		checkError(errors.New("Couldn't parse torrent file"))
	}
	infoDict, ok := metaMap["info"]
	if !ok {
		checkError(errors.New("Unable to locate info dict in torrent file"))
	}

	var b bytes.Buffer
	err = bencode.Marshal(&b, infoDict)
	checkError(err)

	metaInfo := &MetaInfo{}
	file.Seek(0, 0)
	err = bencode.Unmarshal(file, &metaInfo)
	checkError(err)

	return metaInfo
}
