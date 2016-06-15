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
	"io"
	"log"
)

func ParseDict(r *bufio.Reader) map[string]interface{} {
	dict := make(map[string]interface{})

	// dictionaries start with 'd' in bencoding, and a metainfo file is a
	// dictionary
	d_byte, err := r.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	if d_byte != 'd' {
		log.Fatalln("unexpected byte where 'd' for dictionary was expected:", d_byte)
	}

	for {
		byt, err := r.ReadByte()
		if err != nil {
			log.Fatalln(err)
		}
		if byt < 48 || byt > 57 {
			if byt == 'e' {
				break
			} else {
				log.Fatalln("unexpected byte in dictionary metadata:", byt)
			}
		}
		r.UnreadByte()

		// parse key
		key := ParseString(r)

		// parse value
		peek_byte, err := r.Peek(1)
		if err != nil {
			log.Fatalln(err)
		}
		if peek_byte[0] > 47 && peek_byte[0] < 58 {
			dict[key] = ParseString(r)
		} else if peek_byte[0] == 'i' {
			dict[key] = ParseInt(r)
		} else if peek_byte[0] == 'l' {
			dict[key] = ParseList(r)
		} else if peek_byte[0] == 'd' {
			dict[key] = ParseDict(r)
		} else {
			log.Fatalln("unexpected byte in dictionary value metadata:", peek_byte[0])
		}
	}

	return dict
}

func ParseInt(r *bufio.Reader) int64 {
	byt, err := r.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	if byt != 'i' {
		log.Fatalln("unexpected byte where 'i' for integer was expected:", byt)
	}

	var i int64
	// XXX: do we need to check for 0 bytes read?
	_, err = fmt.Fscanf(r, "%de", &i)
	if err != nil {
		log.Fatalln(err)
	}

	return i
}

func ParseList(r *bufio.Reader) []interface{} {
	l := make([]interface{}, 0, 0)

	// dictionaries start with 'd' in bencoding, and a metainfo file is a
	// dictionary
	l_byte, err := r.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	if l_byte != 'l' {
		log.Fatalln("unexpected byte where 'l' for list was expected:", l_byte)
	}

	for {
		byt, err := r.ReadByte()
		if err != nil {
			log.Fatalln(err)
		}
		if byt == 'e' {
			break
		}
		r.UnreadByte()

		// parse value
		peek_byte, err := r.Peek(1)
		if err != nil {
			log.Fatalln(err)
		}
		if peek_byte[0] > 47 && peek_byte[0] < 58 {
			l = append(l, ParseString(r))
		} else if peek_byte[0] == 'i' {
			l = append(l, ParseInt(r))
		} else if peek_byte[0] == 'l' {
			l = append(l, ParseList(r))
		} else if peek_byte[0] == 'd' {
			l = append(l, ParseDict(r))
		} else {
			log.Fatalln("unexpected byte in list value metadata:", peek_byte[0])
		}
	}

	return l
}

func ParseString(r *bufio.Reader) string {
	var sz uint64
	// XXX: do we need to check for 0 bytes read?
	_, err := fmt.Fscanf(r, "%d", &sz)
	if err != nil {
		log.Fatalln(err)
	}

	readSemicolon(r)

	b := make([]byte, sz, sz)

	_, err = io.ReadFull(r, b)
	if err != nil {
		log.Fatalln(err)
	}

	return string(b)
}

func readSemicolon(r *bufio.Reader) {
	byt, err := r.ReadByte()
	if err != nil {
		log.Fatalln(err)
	}
	if byt != ':' {
		log.Fatalf("expected semicolon, encountered %b", byt)
	}
}