package utils

import (
	"chainmaker_web/src/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/test-go/testify/assert"
)

// ActionTest = "/test-action"
var ActionTest = "/test-action"

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")

	// 运行其他测试
	os.Exit(m.Run())
}

func TestGetRelayCrossChainHttpResp1(t *testing.T) {
	params := map[string]string{"key": "value"}
	expectedResponse := []byte(`{"status":"success"}`)

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqParams map[string]string
		err := json.NewDecoder(r.Body).Decode(&reqParams)
		assert.NoError(t, err)
		assert.Equal(t, params, reqParams)

		w.WriteHeader(http.StatusOK)
		w.Write(expectedResponse)
	}))
	defer ts.Close()

	// Override the global config for testing
	config.GlobalConfig.WebConf.RelayCrossChainUrl = ts.URL

	resp, err := GetRelayCrossChainHttpResp(params, ActionTest)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestSendHttpPostRequest(t *testing.T) {
	params := map[string]string{"key": "value"}
	expectedResponse := []byte(`{"status":"success"}`)

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var reqParams map[string]string
		err := json.NewDecoder(r.Body).Decode(&reqParams)
		assert.NoError(t, err)
		assert.Equal(t, params, reqParams)

		w.WriteHeader(http.StatusOK)
		w.Write(expectedResponse)
	}))
	defer ts.Close()

	resp, err := SendHttpPostRequest(ts.URL+ActionTest, ActionTest, params)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestSendHttpGetRequest(t *testing.T) {
	params := map[string]string{"key": "value"}
	expectedResponse := []byte(`{"status":"success"}`)

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		assert.Equal(t, params["key"], query.Get("key"))

		w.WriteHeader(http.StatusOK)
		w.Write(expectedResponse)
	}))
	defer ts.Close()

	resp, err := SendHttpGetRequest(ts.URL, params)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestSendHttpPostRequestNew(t *testing.T) {
	params := map[string]string{"key": "value"}
	expectedResponse := []byte(`{"status":"success"}`)

	// Create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		assert.NoError(t, err)

		assert.Equal(t, params["key"], r.FormValue("key"))

		w.WriteHeader(http.StatusOK)
		w.Write(expectedResponse)
	}))
	defer ts.Close()

	resp, err := SendHttpPostRequestNew(ts.URL, params, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, resp)
}

func TestGetCrossGatewayInfo(t *testing.T) {
	_, _ = GetCrossGatewayInfo(123)
}

func TestGetCrossSubChainInfo(t *testing.T) {
	_, _ = GetCrossSubChainInfo("123")
}
