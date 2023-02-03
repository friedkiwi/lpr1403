package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/klauspost/compress/zstd"
	"io"
	"net/http"
	"os"
)

var serverUrl = flag.String("endpoint", "https://1403.bitnet.systems/print", "API endpoint to submit print job to")
var authToken = flag.String("auth-token", "", "Authorisation token to use to submit print job")

func main() {
	flag.Parse()

	in := bufio.NewReader(os.Stdin)
	rawBuffer := bytes.Buffer{}
	encodedBuffer, _ := zstd.NewWriter(&rawBuffer)
	printOutput := bufio.NewWriter(encodedBuffer)
	var lineCount, pageCount, charCount = 0, 0, 0

	if len(*authToken) < 20 {
		fmt.Printf("Invalid authentication token.\n")
		return
	}

	// loop until EOF
	for {
		// read byte from stdin
		c, err := in.ReadByte()
		if err == io.EOF {
			printOutput.WriteByte(0x0a)
			printOutput.WriteString("J:lpr1403\n")
			printOutput.Flush()
			encodedBuffer.Close()
			break
		}
		if charCount == 0 {
			printOutput.WriteString("L:")
		}

		switch c {
		case 0x0c:
			// form feed
			printOutput.WriteByte(0x0a)
			printOutput.WriteString("P:")
			pageCount++
		case 0x0a:
			// newline
			printOutput.WriteByte(0x0a)
			printOutput.WriteString("L:")
			lineCount++
		default:
			printOutput.WriteByte(c)
		}
		charCount++
	}

	_ = os.WriteFile("debug.zstd", rawBuffer.Bytes(), 0666)

	// post request to print server
	req, err := http.NewRequest(http.MethodPost, (*serverUrl)+"?profile=default-green", &rawBuffer)
	if err != nil {
		fmt.Printf("Error: unable to create HTTP request: %v\n", err)
		return
	}

	req.Header.Set("Content-Encoding", "zstd")
	req.Header.Set("Content-Type", "text/x-print-job")
	req.Header.Set("Authorization", "Bearer "+*authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: unable to execute HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	defer io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Spooled %d lines and %d pages for a total of %d characters.\n", lineCount, pageCount, charCount)
	} else {
		fmt.Printf("ERROR: print API response status: %s\n",
			resp.Status)
	}
}
