package main

import (
	"flag"
	"fmt"
	"os"
)

type CLFlags struct {
	Name         string
	Announce     string
	AnnounceList string
	Comment      string
	CreatedBy    string
	Encoding     string
}

func main() {

	flag.Parse()
	args := flag.Args()
	fmt.Fprintf(os.Stderr, "\n")

	if len(args) > 0 {
		switch args[0] {
		case "maketorrent":
			name := flag.String("name", "", "Tracker name, defaults to [file]+.torrent")
			announce := flag.String("announce", "udp://tracker.publicbt.com:80", "Tracker url")
			announceList := flag.String("announceList", "", "Comma seperated tracker urls")
			comment := flag.String("comment", "", "Optional comment")
			createdBy := flag.String("createdBy", "Tulva", "Author")
			encoding := flag.String("encoding", "UTF-8", "Encoding")

			flag.Parse()
			args = flag.Args()

			if len(args) == 1 {
				fmt.Fprintf(os.Stderr, "Usage: maketorrent [options] File\n\n")
				flag.PrintDefaults()
			} else if len(args) == 2 {
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
			}
		default:
			fmt.Fprintf(os.Stderr, "Available commands:\n\n maketorrent")
		}
	} else {
		fmt.Fprintf(os.Stderr, "Available commands:\n\n maketorrent")
	}
	fmt.Fprintf(os.Stderr, "\n\n")
}
