package tacticchecker

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleTacticSinglePixelSuccess(t *testing.T) {

	testData := []string{
		"", "1", "", "", "", "", "", "", "[\"test url\"]", "", "",
	}

	mockClient := &HTTPClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
			}, nil
		},
	}

	tacticChecker := New(mockClient)

	tacticChecker.CheckTactic(testData)

	tacticChecker.Wg.Wait()

	assert.Equal(t, 1, tacticChecker.SuccessCount)
	assert.Equal(t, 0, tacticChecker.FailureCount)
	assert.Equal(t, 0, len(tacticChecker.FailedPixels))
	assert.Equal(t, 1, len(mockClient.GetCalls()))

}

func TestSingleTacticMultiPixelSuccess(t *testing.T) {

	testData := []string{
		"", "1", "", "", "", "", "", "", "[\"test url\",\"another test url\"]", "", "",
	}

	mockClient := &HTTPClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
			}, nil
		},
	}

	tacticChecker := New(mockClient)

	tacticChecker.CheckTactic(testData)

	tacticChecker.Wg.Wait()

	assert.Equal(t, 2, tacticChecker.SuccessCount)
	assert.Equal(t, 0, tacticChecker.FailureCount)
	assert.Equal(t, 0, len(tacticChecker.FailedPixels))
	assert.Equal(t, 2, len(mockClient.GetCalls()))

}

func TestSingleTacticSinglePixelFail(t *testing.T) {

	testData := []string{
		"", "1", "", "", "", "", "", "", "[\"test url\"]", "", "",
	}

	mockClient := &HTTPClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
			}, nil

		},
	}

	tacticChecker := New(mockClient)

	tacticChecker.CheckTactic(testData)
	tacticChecker.Wg.Wait()

	assert.Equal(t, 0, tacticChecker.SuccessCount)
	assert.Equal(t, 1, tacticChecker.FailureCount)
	assert.Equal(t, 1, len(tacticChecker.FailedPixels))
	assert.Equal(t, "test url", tacticChecker.FailedPixels["1"][0])
	assert.Equal(t, 1, len(mockClient.GetCalls()))

}

func TestSingleTacticMultiPixelFail(t *testing.T) {

	testData := []string{
		"", "1", "", "", "", "", "", "", "[\"test url\",\"another test url\"]", "", "",
	}

	mockClient := &HTTPClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
			}, nil

		},
	}

	tacticChecker := New(mockClient)

	tacticChecker.CheckTactic(testData)
	tacticChecker.Wg.Wait()

	assert.Equal(t, 0, tacticChecker.SuccessCount)
	assert.Equal(t, 2, tacticChecker.FailureCount)
	assert.Equal(t, 1, len(tacticChecker.FailedPixels))

	assert.Equal(t, 2, len(tacticChecker.FailedPixels["1"]))
	for _, val := range tacticChecker.FailedPixels["1"] {
		assert.True(t, val == "test url" || val == "another test url")
	}

	//order is non-deterministic, so we can't do this with multiple items
	//assert.Equal(t, "test url", tacticChecker.FailedPixels["1"][0])
	//assert.Equal(t, "test url", tacticChecker.FailedPixels["1"][1])
	assert.Equal(t, 2, len(mockClient.GetCalls()))

}

// func multiple_tactic_test(t *testing.T) {

// 	testData := []string{
// 		"", "1", "", "", "", "", "", "", "test url", "", "",
// 	}

// 	mockClient := &HTTPClientMock{
// 		GetFunc: func(url string) (*http.Response, error) {

// 			if url != "fail" {
// 				return &http.Response{
// 					StatusCode: 200,
// 				}, nil
// 			} else {
// 				return &http.Response{
// 					StatusCode: 400,
// 				}, nil
// 			}
// 		},
// 	}

// 	tacticChecker := New(mockClient)

// 	tacticChecker.CheckTactic(testData)

// 	assert.Equal(t, 1, tacticChecker.SuccessCount)
// 	assert.Equal(t, 1, len(mockClient.GetCalls()))

// }
