package server

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
)

type Parser struct {
	r *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		r: bufio.NewReader(r),
	}
}

func (p *Parser) ParseRequest() (*types.Request, error) {
	rawHeaders, err := readUntilDelimiter(p.r, []byte("\r\n\r\n"))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(rawHeaders), "\r\n")
	parts := strings.Split(lines[0], " ")
	if len(parts) < 2 {
		return nil, errors.New("invalid request line")
	}

	req := &types.Request{
		Method:  parts[0],
		Path:    parts[1],
		Headers: make(map[string]string),
	}

	for _, line := range lines[1:] {
		if line == "" {
			break
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) == 2 {
			req.Headers[strings.TrimSpace(kv[0])] =
				strings.TrimSpace(kv[1])
		}
	}

	// Body only if Content-Length exists
	if cl, ok := req.Headers["Content-Length"]; ok {
		n, err := strconv.Atoi(cl)
		if err != nil {
			return nil, err
		}
		req.Body = make([]byte, n)
		if _, err := io.ReadFull(p.r, req.Body); err != nil {
			return nil, err
		}
	}

	return req, nil
}
