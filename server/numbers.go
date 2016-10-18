package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	Timeout = 500

	ErrURLQueryParamNotCorrect = "Query parameter not correct"
	ErrURLQueryParamAmbiguous  = "Query parameter ambiguous"
)

// Results are integers returned from an endpoint
type Results struct {
	// Numbers stores list of integers
	Numbers []int `json:"numbers"`
	url     string
	err     error
}

// processQuery extracts url from query parameters, sends requests and prepares json response
func processQuery(values map[string][]string) []byte {
	// check for query param
	urls, ok := values["u"]
	if !ok {
		LogError(ErrURLQueryParamNotCorrect)
		return []byte{}
	}

	// query param ambiguous?
	if len(values) > 1 {
		LogWarn(ErrURLQueryParamAmbiguous)
	}

	results := processURLs(urls...)

	respJson, err := json.Marshal(results)
	if err != nil {
		LogError(err.Error())
		return []byte{}
	}

	return respJson
}

// processURLs calls endpoints for list of integers, aggregates, filters and sorts results
func processURLs(urls ...string) Results {

	pipe := make(chan Results, len(urls))

	// to aggregate results
	results := make([]Results, len(urls))

	// main routine to accumulate results from workers and timeout if too long
	var wg sync.WaitGroup
	wg.Add(1)
	go func(pipe chan Results) {
		defer wg.Done()
		// to count work done by workers
		done := 0
		for {
			select {
			case res := <-pipe:
				LogDebug(fmt.Sprintf("Read from endpoint %s", res.url))
				results[done] = res
				// work done, increase
				done++
			case <-time.After(time.Millisecond * Timeout):
				// nil channel to stop sending
				pipe = nil
				LogWarn("Too long! Timing out ...")
				return
			}
			// break loop if all work is done
			if done == len(urls) {
				LogDebug("Successfully reached all endpoints")
				return
			}
		}
	}(pipe)

	// spawn workers
	for _, url := range urls {
		go getURL(url, pipe)
	}

	// wait until main routine finishes
	wg.Wait()

	// aggregate ints from all results
	var aggregated []int
	for _, res := range results {
		aggregated = append(aggregated, res.Numbers...)
	}

	// create helper map to for easy filtering
	f := make(filter, len(aggregated))
	for _, value := range aggregated {
		f[value] = true
	}

	filtered := f.ints()
	sort.Ints(filtered)
	return Results{Numbers: filtered}
}

func getURL(url string, pipe chan Results) {
	resp, err := http.Get(url)
	if err != nil {
		LogError(err.Error())
		return
	}

	if resp.StatusCode != http.StatusOK {
		LogError(fmt.Sprintf("Response status from %s incorrect - %d", url, resp.StatusCode))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogError(fmt.Sprintf("Can't read response body - %s", err.Error()))
		return
	}

	var results Results
	err = json.Unmarshal(body, &results)
	if err != nil {
		LogError(fmt.Sprintf("Can't decode response body - %s", err.Error()))
		return
	}
	results.url = url
	// prevent sending to nil channel (when timeout occurred in the mean time
	if pipe != nil {
		pipe <- results
	}
}

// filter is alias type used for int filtering
type filter map[int]bool

// ints returns integer filter keys
func (f filter) ints() []int {
	keys := make([]int, len(f))
	i := 0
	for key := range f {
		keys[i] = key
		i++
	}
	return keys
}
