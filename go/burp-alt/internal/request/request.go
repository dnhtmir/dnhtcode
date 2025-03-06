package request

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
)

type Connection struct {
	Request  *http.Request
	Response *http.Response
}

func ReadHTTP2FromFile(r io.Reader) error {
	buf := bufio.NewReader(r)
	buf2 := new(bytes.Buffer)

	// Create a Framer to read HTTP/2 frames
	framer := http2.NewFramer(buf2, buf)

	for {
		frame, err := framer.ReadFrame()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading frame:", err)
			return err
		}

		// Process the frame (this is just an example)
		fmt.Printf("Frame: %v\n", frame)
	}
	return nil
}
