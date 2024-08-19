package collect

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	ErrFetchUrl      = fmt.Errorf("fetch url error")
	ErrBadStatusCode = fmt.Errorf("bad status code error")
)

type Fetcher interface {
	Get(url string) ([]byte, error)
}

type BaseFetch struct {
	Timeout time.Duration
}

var _ Fetcher = (*BaseFetch)(nil)

func NewBaseFetch(timeout time.Duration) *BaseFetch {
	return &BaseFetch{
		Timeout: timeout,
	}
}

// Get implements Fetcher.
func (b *BaseFetch) Get(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("get url failed")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, ErrFetchUrl
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("error status code: %d\n", resp.StatusCode)
		return nil, ErrBadStatusCode
	}

	bodyReader := bufio.NewReader(resp.Body)
	enc := DetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, enc.NewDecoder())
	return io.ReadAll(utf8Reader)
}

func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
