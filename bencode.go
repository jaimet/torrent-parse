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
)

func ParseDict(r *bufio.Reader) (error, map[string]interface{}) {
	dict := make(map[string]interface{})

	// dictionaries start with 'd' in bencoding, and a metainfo file is a
	// dictionary
	d_byte, err := r.ReadByte()
	if err != nil {
		return err, nil
	}
	if d_byte != 'd' {
		return fmt.Errorf("unexpected byte where 'd' for dictionary was expected:", d_byte), nil
	}

	for {
		byt, err := r.ReadByte()
		if err != nil {
			return err, nil
		}
		if byt < 48 || byt > 57 {
			if byt == 'e' {
				break
			} else {
				return fmt.Errorf("unexpected byte in dictionary metadata:", byt), nil
			}
		}
		r.UnreadByte()

		// parse key
		err, key := ParseString(r)
		if err != nil {
			return err, nil
		}

		// parse value
		peek_byte, err := r.Peek(1)
		if err != nil {
			return err, nil
		}
		if peek_byte[0] > 47 && peek_byte[0] < 58 {
			err, dict[key] = ParseString(r)
		} else if peek_byte[0] == 'i' {
			err, dict[key] = ParseInt(r)
		} else if peek_byte[0] == 'l' {
			err, dict[key] = ParseList(r)
		} else if peek_byte[0] == 'd' {
			err, dict[key] = ParseDict(r)
		} else {
			err = fmt.Errorf("unexpected byte in dictionary value metadata:", peek_byte[0])
		}

		if err != nil {
			return err, nil
		}
	}

	return nil, dict
}

func ParseInt(r *bufio.Reader) (error, int64) {
	byt, err := r.ReadByte()
	if err != nil {
		return err, 0
	}
	if byt != 'i' {
		return fmt.Errorf("unexpected byte where 'i' for integer was expected:", byt), 0
	}

	var i int64
	// XXX: do we need to check for 0 bytes read?
	_, err = fmt.Fscanf(r, "%de", &i)
	if err != nil {
		return err, 0
	}

	return nil, i
}

func ParseList(r *bufio.Reader) (error, []interface{}) {
	l := make([]interface{}, 0, 0)

	// dictionaries start with 'd' in bencoding, and a metainfo file is a
	// dictionary
	l_byte, err := r.ReadByte()
	if err != nil {
		return err, nil
	}
	if l_byte != 'l' {
		return fmt.Errorf("unexpected byte where 'l' for list was expected:", l_byte), nil
	}

	for {
		byt, err := r.ReadByte()
		if err != nil {
			return err, nil
		}
		if byt == 'e' {
			break
		}
		r.UnreadByte()

		// parse value
		peek_byte, err := r.Peek(1)
		if err != nil {
			return err, nil
		}
		if peek_byte[0] > 47 && peek_byte[0] < 58 {
			err, s := ParseString(r)
			if err != nil {
				return err, nil
			}
			l = append(l, s)
		} else if peek_byte[0] == 'i' {
			err, i := ParseInt(r)
			if err != nil {
				return err, nil
			}
			l = append(l, i)
		} else if peek_byte[0] == 'l' {
			err, li := ParseList(r)
			if err != nil {
				return err, nil
			}
			l = append(l, li)
		} else if peek_byte[0] == 'd' {
			err, d := ParseDict(r)
			if err != nil {
				return err, nil
			}
			l = append(l, d)
		} else {
			return fmt.Errorf("unexpected byte in list value metadata:", peek_byte[0]), nil
		}
	}

	return nil, l
}

func ParseString(r *bufio.Reader) (error, string) {
	var sz uint64
	// XXX: do we need to check for 0 bytes read?
	_, err := fmt.Fscanf(r, "%d", &sz)
	if err != nil {
		return err, ""
	}

	// jump over semicolon
	byt, err := r.ReadByte()
	if err != nil {
		return err, ""
	}
	if byt != ':' {
		return fmt.Errorf("expected semicolon, encountered %b", byt), ""
	}

	b := make([]byte, sz, sz)

	_, err = io.ReadFull(r, b)
	if err != nil {
		return err, ""
	}

	return nil, string(b)
}
