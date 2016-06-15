/*
 * Copyright (c) 2016 Michael McConville <mmcco@mykolab.com>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var metainfo *bufio.Reader

	switch len(os.Args) {
	case 1:
		metainfo = bufio.NewReader(os.Stdin)
	case 2:
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}
		metainfo = bufio.NewReader(file)
	default:
		log.Fatalln("torrent-parse [metainfo-file]")
	}

	err, d := ParseDict(metainfo)
	if err != nil {
		log.Fatalln(err)
	}
	pretty_print(d)
}

func pretty_print(metainfo map[string]interface{}) {
	for k, v := range metainfo {
		switch k {
		case "announce":
			switch s := v.(type) {
			case string:
				fmt.Printf("\ntracker URL:\t\t%s\n", s)
			default:
				fmt.Println("\ntracker URL:\t\t[unexpected value type]")
			}
		case "info":
			d, ok := v.(map[string]interface{})
			if !ok {
				log.Fatalln("'info' value in metainfo dictionary has unexpected type")
			}

			for ki, vi := range d {
				if ki == "length" {
					switch flen := vi.(type) {
					case int64:
						fmt.Printf("\nfile length:\t\t%d\n", flen)
					default:
						fmt.Println("\nfile length:\t\t[unexpected value type]")
					}
				} else if ki == "files" {
					files, ok := vi.([]interface{})
					if !ok {
						log.Fatalln("'files' value in metainfo dictionary has unexpected type")
					}
					for _, file := range files {
						f, ok := file.(map[string]interface{})
						if !ok {
							log.Fatalln("item in 'files' list has unexpected type")
						}
						fpath, ok := f["path"].([]interface{})
						if !ok {
							log.Fatalln("file in 'files' list has a path of unexpected type")
						}
						// hack to avoid converting and using strings.Join
						for i, dir := range fpath {
							if i == 0 {
								fmt.Print("\n\t")
							} else {
								fmt.Print("/")
							}
							fmt.Print(dir)
						}
						fmt.Println()
						fmt.Printf("\tfile size:\t\t%d\n", f["length"])
					}
				}
			}
		case "created by":
			fmt.Printf("\ncreated with:\t\t%s\n", v)
		case "creation date":
			date, ok := v.(int64)
			if !ok {
				log.Fatalln("torrent creation date has unexpected type")
			}
			fmt.Printf("\ncreation date:\t\t%s\n", time.Unix(date, 0).String())
		case "comment":
			fmt.Printf("\ncomment:\t\t%s\n", v)
		case "encoding":
			fmt.Printf("\nencoding:\t\t%s\n", v)
		case "announce-list":
			fmt.Printf("\nannounce list:\t\t%s\n", v)
		case "url-list":
			fmt.Printf("\nURL list:\t\t%s\n", v)
		case "errors":
			fmt.Printf("\nerrors:\t\t%s\n", v)
		case "err_callback":
			fmt.Printf("\nerror callback:\t\t%s\n", v)
		case "log_callback":
			fmt.Printf("\nlog callback:\t\t%s\n", v)
		case "httpseeds":
			fmt.Printf("\nHTTP seeds:\t\t%s\n", v)
		default:
			fmt.Fprintln(os.Stderr, "\nencountered additional top-level dictionary key:", k)
		}
	}
	fmt.Println()
}
