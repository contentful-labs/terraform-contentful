package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSpace(t *testing.T) {
	cmaToken := "test-token"
	// TODO: If only one org, create without requiring orgID. Two tests
	organizationID := "test-organization"
	spaceName := "test-name"
	defaultLocale := "en"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "POST")
		assert.Equal(t, req.RequestURI, "/spaces/")
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
		assert.Equal(t, organizationID, req.Header.Get("X-Contentful-Organization"))

		var payload map[string]string
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&payload))
		assert.Equal(t, spaceName, payload["name"])
		assert.Equal(t, defaultLocale, payload["defaultLocale"])

		w.WriteHeader(201)
		err := json.NewEncoder(w).Encode(&spaceData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	_, err := createSpace(cmaToken, organizationID, spaceName, defaultLocale)
	assert.Nil(t, err)
}

func TestUpdateSpace(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"
	version := 1
	newSpaceName := "test-name"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "PUT")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))
		assert.Equal(t, "application/vnd.contentful.management.v1+json", req.Header.Get("Content-Type"))
		assert.Equal(t, strconv.Itoa(version), req.Header.Get("X-Contentful-Version"))

		var payload map[string]string
		assert.Nil(t, json.NewDecoder(req.Body).Decode(&payload))
		assert.Equal(t, newSpaceName, payload["name"])

		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(&spaceData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	_, err := updateSpace(cmaToken, spaceID, version, newSpaceName)
	assert.Nil(t, err)
}

func TestReadSpace(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "GET")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))

		w.WriteHeader(200)
		err := json.NewEncoder(w).Encode(&spaceData{})
		assert.Nil(t, err)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	_, err := readSpace(cmaToken, spaceID)
	assert.Nil(t, err)
}

func TestDeleteSpace(t *testing.T) {
	cmaToken := "test-token"
	spaceID := "test-space"

	serveFn := func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.Method, "DELETE")
		assert.Equal(t, req.RequestURI, "/spaces/"+spaceID)
		assert.Equal(t, "Bearer "+cmaToken, req.Header.Get("Authorization"))

		w.WriteHeader(204)
	}

	server := httptest.NewServer(&testHandler{handleFunc: serveFn})
	defer server.Close()
	baseURL = server.URL

	err := deleteSpace(cmaToken, spaceID)
	assert.Nil(t, err)
}

type httpClientMock struct {
	passedInRequest  *http.Request
	returnedResponse *http.Response
	returnedError    error
}

func (c *httpClientMock) Do(req *http.Request) (*http.Response, error) {
	c.passedInRequest = req
	return c.returnedResponse, c.returnedError
}
