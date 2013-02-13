package main

import (
	"io"
	"os"
)

// Encodes the given byte slice by returning a separate byte slice in which
// every byte of the output slice is 13 higher than the corresponding byte
// in the input slice, wrapping around when greater than 127
func Encode(in []byte) []byte {
	out := make([]byte, len(in))

	for i := 0; i < len(in); i++ {
		out[i] = (in[i] + 13) % 128
	}

	return out
}

// Decodes an input byte slice that was previously encoded with Encode
func Decode(in []byte) []byte {
	out := make([]byte, len(in))

	for i := 0; i < len(in); i++ {
		b := in[i]
		if b < 13 {
			b += 128
		}
		out[i] = (b - 13) % 128
	}

	return out
}

// Given a plaintext serialization of a speaker map, newline-delimited of the
// form key,value, creates an encoded serialization
func encodeSpeakers(speakers *os.File) []byte {
	contents := make([]byte, 0)
	var err error
	for err == nil { 
		buf := make([]byte, 1)
		_, err = io.ReadFull(speakers, buf)
		if err == nil {
			contents = append(contents, buf[:len(buf)]...)
		}
	}

	return Encode(contents)
}

// Decodes a serialization previously returned by encodeSpeakers, to be parsed
// and reconstructed into a map
func decodeSpeakers(encodedSpeakers *os.File) []byte {
	contents := make([]byte, 0)
	var err error
	for err == nil { 
		buf := make([]byte, 1)
		_, err = io.ReadFull(encodedSpeakers, buf)
		if err == nil {
			contents = append(contents, buf[:len(buf)]...)
		}
	}

	return Decode(contents)
}
