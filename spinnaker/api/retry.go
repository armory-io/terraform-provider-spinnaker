package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/cenkalti/backoff"
)

const (
	// maxRetries defines the number of times an operation should be retried in
	// the worst case.
	maxRetries = 5

	// retryInterval specified the time to wait between two retries.
	retryInterval = 2 * time.Second
)

// operation is a func that performs an API call and returns the response json
// data as map, the http response and an error.
type operation func() (map[string]interface{}, *http.Response, error)

// retry retries an operation using constant backoff.
func retry(fn operation) (data map[string]interface{}, resp *http.Response, err error) {
	err = backoff.Retry(func() error {
		data, resp, err = fn()

		// Check the response code. We retry on 500-range responses to allow
		// the server time to recover, as 500's are typically not permanent
		// errors and may relate to outages on the server side. This will catch
		// invalid response codes as well, like 0 and 999.
		if resp != nil && (resp.StatusCode == 0 || resp.StatusCode >= 500) {
			return errors.New(http.StatusText(resp.StatusCode))
		}

		return err
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(retryInterval), maxRetries))

	return
}
