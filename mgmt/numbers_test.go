package mgmt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

const (
	port1 = "12345"
	port2 = "12346"
)

// TestGetURLsPositive is testing happy path for retrieving  numbers from endpoint
func TestGetURLsPositive(t *testing.T) {

	expected := []int{1, 2, 3}
	endpoint := "/foo"
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port1),
		Path:   endpoint,
	}

	go mockServer(endpoint, port1, expected)

	pipe := make(chan Results, 1)
	getURL(u.String(), pipe)

LOOP:
	for {
		select {
		case res := <-pipe:
			if len(res.Numbers) != len(expected) {
				t.Fail()
				break LOOP
			}

			for i, num := range res.Numbers {
				if num != expected[i] {
					t.Fail()
				}

			}
			break LOOP
		case <-time.After(time.Millisecond * 500):
			t.Fail()
			break LOOP
		}
	}
}

// TestGetURLsNegative is testing failed attempt to retrieve numbers from endpoint
func TestGetURLsNegative(t *testing.T) {

	expected := []int{1, 2, 3}
	endpoint := "/foo"
	bad := "/fool"
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port1),
		Path:   bad,
	}

	go mockServer(endpoint, port1, expected)

	pipe := make(chan Results, 1)
	getURL(u.String(), pipe)

LOOP:
	for {
		select {
		case <-pipe:
			t.Fail()
			break LOOP
		case <-time.After(time.Millisecond * 500):
			break LOOP
		}

	}
}

// TestProcessURLsPositive is testing happy path for processing numbers from given urls (single url)
func TestProcessURLsPositive(t *testing.T) {
	expected := []int{1, 2, 3}
	numbers := []int{3, 1, 2, 2}
	endpoint := "/foo"
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port1),
		Path:   endpoint,
	}

	go mockServer(endpoint, port1, numbers)

	results := processURLs(u.String())

	if len(results.Numbers) != len(expected) {
		t.Fail()
		return
	}

	for i, exp := range expected {
		if exp != results.Numbers[i] {
			t.Fail()
		}
	}
}

// TestProcessQueryPositive is testing happy path for processing query (single url)
func TestProcessQueryPositive(t *testing.T) {
	numbers := []int{1, 2, 3}
	endpoint := "/foo"
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port1),
		Path:   endpoint,
	}

	go mockServer(endpoint, port1, numbers)
	expected, _ := json.Marshal(Results{Numbers: numbers})
	resp := processQuery(map[string][]string{"u": []string{u.String()}})

	if len(resp) != len(expected) {
		t.Fail()
		return
	}

	for i, exp := range expected {
		if exp != resp[i] {
			t.Fail()
		}
	}
}

// TestProcessQueryNegative is testing failed attempt to process query (bad parameter)
func TestProcessQueryNegative(t *testing.T) {
	numbers := []int{1, 2, 3}
	endpoint := "/foo"
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port1),
		Path:   endpoint,
	}

	go mockServer(endpoint, port1, numbers)
	expected, _ := json.Marshal(Results{Numbers: numbers})
	resp := processQuery(map[string][]string{"uu": []string{u.String()}})

	if len(resp) == len(expected) {
		t.Fail()
	}

	if len(resp) != 0 {
		t.Fail()
	}
}

// TestProcessQueryPositive is testing happy path for processing query (two urls)
func TestProcessQueryPositiveTwo(t *testing.T) {
	numbers1 := []int{1, 2, 3}
	endpoint1 := "/foo"
	u1 := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port1),
		Path:   endpoint1,
	}
	numbers2 := []int{3, 4, 1}
	endpoint2 := "/foo"
	u2 := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("localhost:%s", port2),
		Path:   endpoint2,
	}
	combined := []int{1, 2, 3, 4}

	go mockServer(endpoint1, port1, numbers1)
	go mockServer(endpoint2, port2, numbers2)

	expected, _ := json.Marshal(Results{Numbers: combined})
	resp := processQuery(map[string][]string{"u": []string{u1.String(), u2.String()}})

	if len(resp) != len(expected) {
		t.Fail()
		return
	}

	for i, exp := range expected {
		if exp != resp[i] {
			t.Fail()
		}
	}
}

func mockServer(endpoint string, port string, expected []int) {
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string][]int{"numbers": expected}
		jsonResp, _ := json.Marshal(resp)
		fmt.Println(string(jsonResp))
		w.Write(jsonResp)
	})
	http.ListenAndServe(":"+port, mux)
}
