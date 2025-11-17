package utils

import (
	"bytes"
	"chainmaker_web/src/config"
	loggers "chainmaker_web/src/logger"
	"chainmaker_web/src/models/relayCrossChain"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

var log = loggers.GetLogger(loggers.MODULE_WEB)

func GetRelayCrossChainHttpResp(params interface{}, action string) ([]byte, error) {
	crossChainUrl := config.GlobalConfig.WebConf.RelayCrossChainUrl
	url := crossChainUrl + action
	jsonByte, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if resp.Body != nil {
			err = resp.Body.Close()
			if err != nil {
				return
			}
		}
	}()

	if resp.StatusCode != 200 {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	log.Infof("【http service】get relay cross log, params:%v, url:%v, respJson:%v, err:%v",
		string(jsonByte), url, string(body), err)
	return body, err
}

// SendHttpPostRequest sends an HTTP POST request to the specified URL with the given action and parameters.
// It returns the response body as a byte slice and any error encountered.
func SendHttpPostRequest(url string, action string, params interface{}) ([]byte, error) {
	// Marshal the parameters to JSON
	jsonByte, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal params: %v", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer func() {
		if resp.Body != nil {
			if errClose := resp.Body.Close(); errClose != nil {
				log.Errorf("failed to close response body: %v", errClose)
			}
		}
	}()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the request and response
	log.Infof("HTTP request sent. URL: %s, Action: %s, Params: %s, Response: %s",
		url, action, string(jsonByte), string(body))

	return body, nil
}

// SendHttpGetRequest sends an HTTP GET request to the specified URL with the given parameters.
// It returns the response body as a byte slice and any error encountered.
func SendHttpGetRequest(baseURL string, params map[string]string) ([]byte, error) {
	// Parse the base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// Add query parameters to the URL
	query := parsedURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer func() {
		if resp.Body != nil {
			if errClose := resp.Body.Close(); errClose != nil {
				log.Errorf("failed to close response body: %v", errClose)
			}
		}
	}()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the request and response
	log.Infof("HTTP request sent. URL: %s, Response: %s", parsedURL.String(), string(body))
	return body, nil
}

// GetCrossGatewayInfo
//
//	@Description: 根据gatewayId获取子链列表
//	@param gatewayId 网关id
//	@return *relayCrossChain.GatewayInfoData 子链列表
//	@return error
func GetCrossGatewayInfo(gatewayId int64) (*relayCrossChain.GatewayInfoData, error) {
	params := relayCrossChain.GetGatewayInfoReq{
		GatewayId: strconv.FormatInt(gatewayId, 10),
	}
	body, err := GetRelayCrossChainHttpResp(params, CrossGateWayIdUrl)
	if err != nil {
		log.Warnf("cross chain GetCrossGatewayInfo err:%v", err)
		return nil, err
	}
	var respJson relayCrossChain.GetGatewayInfoResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Warnf("cross chain GetCrossGatewayInfo err:%v", err)
		return nil, err
	}

	log.Infof("cross chain GetCrossGatewayInfo GatewayId:%v, result:%v", params.GatewayId, string(body))
	return respJson.Data, nil
}

// GetCrossSubChainInfo
//
//	@Description: 根据子链id获取子链信息
//	@param subChainId
//	@return *relayCrossChain.SubChainInfoData
//	@return error
func GetCrossSubChainInfo(subChainId string) (*relayCrossChain.SubChainInfoData, error) {
	params := relayCrossChain.GetSubChainInfoReq{
		SubChainId: subChainId,
	}
	body, err := GetRelayCrossChainHttpResp(params, CrossSubChainInfoUrl)
	if err != nil {
		return nil, err
	}
	var respJson relayCrossChain.GetSubChainInfoResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		return nil, err
	}

	return respJson.Data, nil
}

func SendHttpPostRequestNew(url string, params map[string]string, files map[string]*multipart.FileHeader) (
	[]byte, error) {
	// 创建一个新的表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加表单字段
	for key, value := range params {
		_ = writer.WriteField(key, value)
	}

	// 添加文件
	for filename, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()
		part, err := writer.CreateFormFile(filename, fileHeader.Filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %v", err)
		}
		_, _ = io.Copy(part, file)
	}

	// 关闭表单
	writer.Close()

	// 创建一个新的请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 记录请求和响应
	log.Infof("HTTP request sent. URL: %s, Params: %v, Response: %s",
		url, params, string(bodyBytes))

	return bodyBytes, nil
}
