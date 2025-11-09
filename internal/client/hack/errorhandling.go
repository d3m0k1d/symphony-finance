package hack

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// format must contain one %w
func fmtHttpError(format string, resp *http.Response) (err error) {
	bits, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error while handling http error: %w", err)
	}
	defer func() {

		resp.Body = io.NopCloser(bytes.NewReader(bits))
	}()

	return fmt.Errorf(format, fmt.Errorf("status code: %d\nbody:\n%s", resp.StatusCode, bits))
}
