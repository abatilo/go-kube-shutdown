package ready

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPGet returns a Check that performs an HTTP GET request against the
// specified URL. The check fails if the response takes longer than 200ms or
// returns a non-200 status code.
func HTTPGet(url string) Check {
	client := http.Client{
		// Opinionated limit for how slow a server can respond
		Timeout: time.Millisecond * 200,
		// never follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return func() error {
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Received %d", resp.StatusCode)
		}
		return nil
	}
}

// HTTPGetWithTimeout returns a Check that performs an HTTP GET request against
// the specified URL. The check fails if the response takes longer than the
// specified timeout or returns a non-200 status code.
func HTTPGetWithTimeout(url string, timeout time.Duration) Check {
	client := http.Client{
		// Opinionated limit for how slow a server can respond
		Timeout: timeout,
		// never follow redirects
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return func() error {
		resp, err := client.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Received %d", resp.StatusCode)
		}
		return nil
	}
}
