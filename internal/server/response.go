package server

import (
	"fmt"
	"io"

	"github.com/Dishank-Sen/Blockchain-Scratch-Bootstrap/types"
)

func writeResponse(w io.Writer, resp *types.Response) error {
	if resp.Headers == nil {
		resp.Headers = make(map[string]string)
	}
	resp.Headers["Content-Length"] = fmt.Sprintf("%d", len(resp.Body))

	if _, err := fmt.Fprintf(w, "%d %s\r\n", resp.StatusCode, resp.Message); err != nil {
		return err
	}

	for k, v := range resp.Headers {
		if _, err := fmt.Fprintf(w, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}

	_, err := w.Write(resp.Body)
	return err
}
