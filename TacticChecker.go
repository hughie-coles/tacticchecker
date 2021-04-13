package tacticchecker

//go:generate moq -out MockHTTPClient.go . HTTPClient
import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type ConcreteHttpClient struct {
	Client http.Client
}

func (c ConcreteHttpClient) Get(url string) (*http.Response, error) {
	return c.Client.Get(url)
}

type Interface interface {
	checkTactic(tactic []string)
	checkPixel(pixelUrl string) error
}

func New(client HTTPClient) TacticChecker {
	return TacticChecker{
		httpClient:   client,
		FailedPixels: make(map[string][]string, 0),
	}

}

type TacticChecker struct {
	httpClient   HTTPClient
	Wg           sync.WaitGroup
	mutex        sync.Mutex
	SuccessCount int
	FailureCount int
	FailedPixels map[string][]string
}

func (tc *TacticChecker) CheckTactic(tactic []string) {

	impressionPixelJsonColumn := tactic[8]
	var pixelsInTactic []string
	json.Unmarshal([]byte(impressionPixelJsonColumn), &pixelsInTactic)

	for _, pixelUrl := range pixelsInTactic {
		tc.Wg.Add(1)
		go func(pixelUrl string) {
			err := tc.checkPixel(pixelUrl)
			if err != nil {
				tc.mutex.Lock()
				if _, ok := tc.FailedPixels[tactic[1]]; !ok {
					tc.FailedPixels[tactic[1]] = make([]string, 0)
				}
				tc.FailedPixels[tactic[1]] = append(tc.FailedPixels[tactic[1]], pixelUrl)
				tc.FailureCount = tc.FailureCount + 1
				tc.mutex.Unlock()
			} else {
				tc.mutex.Lock()
				tc.SuccessCount = tc.SuccessCount + 1
				tc.mutex.Unlock()
			}
			tc.Wg.Done()
		}(pixelUrl)
	}
}

func (tc *TacticChecker) checkPixel(pixelUrl string) error {
	resp, err := tc.httpClient.Get(pixelUrl)

	// if we error out, mark this as failed.  I'd rather have a false negative
	if err != nil {
		return err
	}

	// 2xx or 3xx == success, anything else is a failure (FYI not strictly the same as the instructions, but I can't see it returning 1xx)
	if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
		return nil
	} else {
		return errors.New("pixel failed")
	}
}
