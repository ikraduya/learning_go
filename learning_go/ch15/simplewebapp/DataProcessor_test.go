package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parser(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected Input
		errMsg   string
	}{
		{"valid", "valid\n+\n1\n2", Input{"valid", "+", 1, 2}, ""},
		{"val1_not_int", "val1_not_int\n+\n1.1\n2", Input{}, "strconv.Atoi: parsing \"1.1\": invalid syntax"},
		{"val2_not_int", "val2_not_int\n+\n1\n2.3", Input{}, "strconv.Atoi: parsing \"2.3\": invalid syntax"},
		{"val1_empty", "val1_empty\n+\n\n2", Input{}, "strconv.Atoi: parsing \"\": invalid syntax"},
		{"val2_empty", "val2_empty\n+\n1\n", Input{}, "strconv.Atoi: parsing \"\": invalid syntax"},
		{"3_lines", "3_lines\n+\n1", Input{}, "entry must be consist of 4 lines"},
		{"2_lines", "2_lines\n+", Input{}, "entry must be consist of 4 lines"},
		{"1_lines", "1_lines", Input{}, "entry must be consist of 4 lines"},
		{"empty", "", Input{}, "entry must be consist of 4 lines"},
		{"5_lines", "5_lines\n+\n3\n4\nbaba", Input{"5_lines", "+", 3, 4}, ""},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			input, err := parser([]byte(d.input))
			if input != d.expected {
				t.Errorf("Expected %v, got %v", d.expected, input)
			}

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.errMsg {
				t.Errorf("Expected error message `%s`, got `%s`", d.errMsg, errMsg)
			}
		})
	}
}

func TestDataProcessor(t *testing.T) {
	data := []struct {
		name     string
		input    string
		expected Result
	}{
		{"plus", "plus\n+\n1\n2", Result{"plus", 3}},
		{"minus", "minus\n-\n1\n2", Result{"minus", -1}},
		{"times", "times\n*\n1\n2", Result{"times", 2}},
		{"division", "division\n/\n1\n2", Result{"division", 0}},
		{"invalid_op", "invalid_op\n^\n1\n2", Result{}},
		{"invalid_data", "invalid_data\n^\n1", Result{}},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			in := make(chan []byte)
			out := make(chan Result)

			go DataProcessor(in, out)
			in <- []byte(d.input)

			if d.name == "invalid_op" || d.name == "invalid_data" {
				close(in)
				return
			}

			result := <-out
			close(in)
			if result != d.expected {
				t.Errorf("Expected %v, got %v", d.expected, result)
			}
		})
	}
}

type MockWriter struct {
	p chan []byte
}

func newMockWriter() *MockWriter {
	return &MockWriter{make(chan []byte, 1)}
}

func (w *MockWriter) Write(p []byte) (n int, err error) {
	w.p <- p
	return 0, nil
}

func TestWriteData(t *testing.T) {
	in := make(chan Result)

	w := newMockWriter()
	go WriteData(in, w)

	in <- Result{"id", 10}
	close(in)
	expected := "id:10\n"
	pOut := string(<-w.p)
	if pOut != expected {
		t.Errorf("Expected `%s`, got `%s`", expected, pOut)
	}
}

type RemoteSolver struct {
	ServerURL string
	Client    *http.Client
}

type CustomErrorReader struct {
}

func (r CustomErrorReader) Read(p []byte) (n int, err error) {
	return len(p), errors.New("Custom Error Reader")
}

func (r *CustomErrorReader) Close() error {
	return nil
}

func TestNewController(t *testing.T) {
	outCh := make(chan []byte, 1)

	handler := NewController(outCh)

	data := []struct {
		name            string
		tcType          string
		dataIn          string
		expectedCode    int
		expectedBodyOut string
	}{
		{"case1", "OK", "case1", http.StatusAccepted, "OK: 1"},
		{"case2", "OK", "case2", http.StatusAccepted, "OK: 2"},
		{"case3", "OK", "case3", http.StatusAccepted, "OK: 3"},
		{"bad", "bad", "", http.StatusBadRequest, "Bad Input"},
		{"case4", "OK", "case4", http.StatusAccepted, "OK: 5"},
		{"busy1", "busy", "busy1", http.StatusServiceUnavailable, "Too Busy: 1"},
		{"busy2", "busy", "busy2", http.StatusServiceUnavailable, "Too Busy: 2"},
		{"case5", "OK", "case5", http.StatusAccepted, "OK: 8"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			var req *http.Request
			if d.tcType == "OK" || d.tcType == "busy" {
				req = httptest.NewRequest(http.MethodGet, "/", bytes.NewBufferString(d.dataIn))
			} else {
				req = httptest.NewRequest(http.MethodGet, "/", CustomErrorReader{})
			}
			rr := httptest.NewRecorder()

			if d.tcType == "busy" {
				outCh <- []byte("junk")
			}

			handler.ServeHTTP(rr, req)

			if rr.Code != d.expectedCode {
				t.Errorf("Expected status code %d, got %d", d.expectedCode, rr.Code)
			}

			contents, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal("Failed reading response body", err)
			}
			if string(contents) != d.expectedBodyOut {
				t.Errorf("Expected response body `%s`, got `%s`", d.expectedBodyOut, string(contents))
			}

			switch d.tcType {
			case "OK":
				outContent := <-outCh
				if string(outContent) != d.dataIn {
					t.Errorf("Expected written data to channel `%s`, got `%s`", d.dataIn, outContent)
				}
			case "busy":
				<-outCh
			}
		})
	}
}

func toData(in Input) []byte {
	return []byte(fmt.Sprintf("%s\n%s\n%d\n%d", in.Id, in.Op, in.Val1, in.Val2))
}

func Fuzz_parser(f *testing.F) {
	testcases := [][]byte{
		[]byte("valid\n+\n1\n2"),
		[]byte("5_lines\n+\n3\n4\nbaba"),
		[]byte(""),
		[]byte("val1_not_int\n+\n1.1\n2"),
		[]byte("val2_empty\n+\n1\n"),
		[]byte("3_lines\n+\n1"),
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		out, err := parser(in)
		if err != nil {
			t.Skip("handled error")
		}
		roundTrip := toData(out)
		out2, err := parser(roundTrip)
		if diff := cmp.Diff(out, out2); diff != "" {
			t.Error(diff)
		}
	})
}
