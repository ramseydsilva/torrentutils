package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	name := flag.String("name", "", "Tracker name, defaults to [file]+.torrent")
	announce := flag.String("announce", "udp://tracker.publicbt.com:80", "Tracker url")
	announceList := flag.String("announceList", "", "Comma seperated tracker urls")
	comment := flag.String("comment", "", "Optional comment")
	createdBy := flag.String("createdBy", "Tulva", "Author")
	encoding := flag.String("encoding", "UTF-8", "Encoding")

	flag.Parse()
	args := flag.Args()

	fmt.Fprintf(os.Stderr, "\n")

	switch {
	case len(args) == 1 && args[0] == "mktorrent":
		fmt.Fprintf(os.Stderr, "Usage: mktorrent [options] File\n\n")
		flag.PrintDefaults()
	case len(args) == 2 && args[0] == "mktorrent":
		clf := &CLFlags{
			Name:         *name,
			Announce:     *announce,
			AnnounceList: *announceList,
			Comment:      *comment,
			CreatedBy:    *createdBy,
			Encoding:     *encoding,
		}
		t := MakeTorrentFile(args[1], clf)
		fmt.Fprintf(os.Stderr, "Made torrent File: %s", t.Name())
	default:
		fmt.Fprintf(os.Stderr, "Available commands:\n\n mktorrent")
	}
	fmt.Fprintf(os.Stderr, "\n\n")
}
