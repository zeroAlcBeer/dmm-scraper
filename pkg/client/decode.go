package client

import (
	"bufio"
	"bytes"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
)

func ToUtf8Encoding(body io.Reader) (r io.Reader, name string, certain bool, err error) {

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return
	}
	e, name, certain, err := DetermineEncodingFromReader(bytes.NewReader(b))
	if err != nil {
		return
	}

	r = transform.NewReader(bytes.NewReader(b), e.NewDecoder())
	return
}

func DetermineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	return
}
