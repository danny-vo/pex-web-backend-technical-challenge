package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockFibSequence struct {
	current  uint64
	next     uint64
	previous uint64
}

func (mfs mockFibSequence) GetCurrent(s *Server) uint64 {
	return mfs.current
}

func (mfs mockFibSequence) GetNext(s *Server) uint64 {
	return mfs.next
}

func (mfs mockFibSequence) GetPrevious(s *Server) uint64 {
	return mfs.previous
}

func TestServer_handleCurrent(t *testing.T) {
	type fields struct {
		mfs fibonacciSequence
	}

	type wants struct {
		contentType string
		payload     string
		statusCode  int
	}
	tests := []struct {
		name   string
		fields fields
		wants  wants
	}{
		{
			name: "happy path",
			fields: fields{
				mfs: mockFibSequence{
					current:  5,
					next:     8,
					previous: 3,
				},
			},
			wants: wants{
				contentType: "application/json",
				payload:     fmt.Sprintf(`{"current": %d}`, 5),
				statusCode:  http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}
			fibSeq = tt.fields.mfs
			req := httptest.NewRequest(http.MethodGet, "http://0.0.0.0:8080/current", nil)
			rw := httptest.NewRecorder()

			server.handleCurrent()(rw, req)
			resp := rw.Result()
			payload, _ := ioutil.ReadAll(resp.Body)

			if !reflect.DeepEqual(tt.wants.statusCode, resp.StatusCode) {
				t.Errorf(
					"Incorrect status code written, wanted: %v but got: %v",
					tt.wants.statusCode, resp.StatusCode,
				)
			}

			if !reflect.DeepEqual(tt.wants.contentType, resp.Header.Get("Content-Type")) {
				t.Errorf(
					"Incorrect content type, wanted: %v but got: %v",
					tt.wants.contentType, resp.Header.Get("Content-Type"),
				)
			}

			if !reflect.DeepEqual(tt.wants.payload, string(payload)) {
				t.Errorf(
					"Incorrect payload received, wanted: %s but got: %s",
					tt.wants.payload, string(payload),
				)
			}
		})
	}
}

func TestServer_handleNext(t *testing.T) {
	type fields struct {
		mfs fibonacciSequence
	}

	type wants struct {
		contentType string
		payload     string
		statusCode  int
	}
	tests := []struct {
		name   string
		fields fields
		wants  wants
	}{
		{
			name: "happy path",
			fields: fields{
				mfs: mockFibSequence{
					current:  5,
					next:     8,
					previous: 3,
				},
			},
			wants: wants{
				contentType: "application/json",
				payload:     fmt.Sprintf(`{"next": %d}`, 8),
				statusCode:  http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}
			fibSeq = tt.fields.mfs
			req := httptest.NewRequest(http.MethodGet, "http://0.0.0.0:8080/current", nil)
			rw := httptest.NewRecorder()

			server.handleNext()(rw, req)
			resp := rw.Result()
			payload, _ := ioutil.ReadAll(resp.Body)

			if !reflect.DeepEqual(tt.wants.statusCode, resp.StatusCode) {
				t.Errorf(
					"Incorrect status code written, wanted: %v but got: %v",
					tt.wants.statusCode, resp.StatusCode,
				)
			}

			if !reflect.DeepEqual(tt.wants.contentType, resp.Header.Get("Content-Type")) {
				t.Errorf(
					"Incorrect content type, wanted: %v but got: %v",
					tt.wants.contentType, resp.Header.Get("Content-Type"),
				)
			}

			if !reflect.DeepEqual(tt.wants.payload, string(payload)) {
				t.Errorf(
					"Incorrect payload received, wanted: %s but got: %s",
					tt.wants.payload, string(payload),
				)
			}
		})
	}
}

func TestServer_handlePrevious(t *testing.T) {
	type fields struct {
		mfs fibonacciSequence
	}

	type wants struct {
		contentType string
		payload     string
		statusCode  int
	}
	tests := []struct {
		name   string
		fields fields
		wants  wants
	}{
		{
			name: "happy path",
			fields: fields{
				mfs: mockFibSequence{
					current:  5,
					next:     8,
					previous: 3,
				},
			},
			wants: wants{
				contentType: "application/json",
				payload:     fmt.Sprintf(`{"previous": %d}`, 3),
				statusCode:  http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}
			fibSeq = tt.fields.mfs
			req := httptest.NewRequest(http.MethodGet, "http://0.0.0.0:8080/current", nil)
			rw := httptest.NewRecorder()

			server.handlePrevious()(rw, req)
			resp := rw.Result()
			payload, _ := ioutil.ReadAll(resp.Body)

			if !reflect.DeepEqual(tt.wants.statusCode, resp.StatusCode) {
				t.Errorf(
					"Incorrect status code written, wanted: %v but got: %v",
					tt.wants.statusCode, resp.StatusCode,
				)
			}

			if !reflect.DeepEqual(tt.wants.contentType, resp.Header.Get("Content-Type")) {
				t.Errorf(
					"Incorrect content type, wanted: %v but got: %v",
					tt.wants.contentType, resp.Header.Get("Content-Type"),
				)
			}

			if !reflect.DeepEqual(tt.wants.payload, string(payload)) {
				t.Errorf(
					"Incorrect payload received, wanted: %s but got: %s",
					tt.wants.payload, string(payload),
				)
			}
		})
	}
}

func TestServer_handleHealth(t *testing.T) {
	type wants struct {
		contentType string
		payload     string
		statusCode  int
	}
	tests := []struct {
		name  string
		wants wants
	}{
		{
			name: "happy path",
			wants: wants{
				contentType: "application/json",
				payload:     fmt.Sprintf(`{"status": "healthy"}`),
				statusCode:  http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &Server{}
			req := httptest.NewRequest(http.MethodGet, "http://0.0.0.0:8080/current", nil)
			rw := httptest.NewRecorder()

			server.handleHealth()(rw, req)
			resp := rw.Result()
			payload, _ := ioutil.ReadAll(resp.Body)

			if !reflect.DeepEqual(tt.wants.statusCode, resp.StatusCode) {
				t.Errorf(
					"Incorrect status code written, wanted: %v but got: %v",
					tt.wants.statusCode, resp.StatusCode,
				)
			}

			if !reflect.DeepEqual(tt.wants.contentType, resp.Header.Get("Content-Type")) {
				t.Errorf(
					"Incorrect content type, wanted: %v but got: %v",
					tt.wants.contentType, resp.Header.Get("Content-Type"),
				)
			}

			if !reflect.DeepEqual(tt.wants.payload, string(payload)) {
				t.Errorf(
					"Incorrect payload received, wanted: %s but got: %s",
					tt.wants.payload, string(payload),
				)
			}
		})
	}
}
