package main

import (
	"bytes"
	"fmt"
	"io"
)

// reader reads from `r` if the stream starts with a given signature.
// otherwise it stops and returns with an error.
type signatureReader struct {
	r         io.Reader // reads from the response.Body (or from any reader)
	signature []byte    // stream should start with this initial signature
}

// Read implements the io.Reader interface.
func (sr *signatureReader) Read(b []byte) (n int, err error) {
	n, err = sr.r.Read(b)

	l := len(sr.signature)
	if l == 0 {
		// simply return if the signature has already been detected.
		return
	}
	// 1) buffer   : **** -> b[:3]            -> ***
	//    signature: ***  -> sr.signature[:3] -> ***
	// 2) buffer   : **   -> b[:2]            -> **
	//    signature: **** -> sr.signature[:2] -> **
	if lb := len(b); lb < l {
		l = lb
	}
	if got, want := b[:l], sr.signature[:l]; !bytes.Equal(got, want) {
		err = fmt.Errorf("signature doesn't match, got: % x, want: % x", got, want)
	}
	// Assuming the `len(b)` is 4.
	// 1st Read(): pr.signature[0:4] -> first part
	// 2nd Read(): pr.signature[0:4] -> second part
	sr.signature = sr.signature[l:]
	return n, err
}

// create a reader for detecting only the png signatures.
func pngReader(r io.Reader) io.Reader {
	return &signatureReader{
		r:         r,
		signature: []byte("\x89PNG\r\n\x1a\n"),
	}
}
