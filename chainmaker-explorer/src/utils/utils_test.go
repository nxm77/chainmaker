package utils

import (
	"chainmaker_web/src/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

var TestCrt = `-----BEGIN CERTIFICATE-----
MIICdjCCAhygAwIBAgIDDnGwMAoGCCqBHM9VAYN1MIGKMQswCQYDVQQGEwJDTjEQ
MA4GA1UECBMHQmVpamluZzEQMA4GA1UEBxMHQmVpamluZzEfMB0GA1UEChMWd3gt
b3JnMS5jaGFpbm1ha2VyLm9yZzESMBAGA1UECxMJcm9vdC1jZXJ0MSIwIAYDVQQD
ExljYS53eC1vcmcxLmNoYWlubWFrZXIub3JnMB4XDTI0MTIxMTA3NTQzM1oXDTI5
MTIxMDA3NTQzM1owgY8xCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAw
DgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3Jn
MQ4wDAYDVQQLEwVhZG1pbjErMCkGA1UEAxMiYWRtaW4xLnNpZ24ud3gtb3JnMS5j
aGFpbm1ha2VyLm9yZzBZMBMGByqGSM49AgEGCCqBHM9VAYItA0IABG7grYTnj027
9Whoq6oHR+s3WzDn++H8kPq+yg7UA4ohDNmagQz7PX3lycVM9onQ+YCBmpw3luX4
8GngCmZaGaijajBoMA4GA1UdDwEB/wQEAwIGwDApBgNVHQ4EIgQgRf1tVCQs2n8m
2KAotbwaXqRUZyd73cvh/+Emctq1z0owKwYDVR0jBCQwIoAgikYZ624olnxP1+++
tJF82ZVA98mNLSimJscsIaVteIcwCgYIKoEcz1UBg3UDSAAwRQIhAIZXm6hB2Khg
DWjo8ufMmIDq/DnSomZGXtjro5JcwHGNAiB/t8izCeezbOwkQF80kT+3mAUATYNL
CujpEVYVNeH92A==
-----END CERTIFICATE-----
`

var MemberInfoTest = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNkakNDQWh5Z0F3SUJBZ0lEQlNvNU1Bb0dDQ3FHU000OUJBTUNNSUdLTVFzd0NRWURWUVFHRXdKRFRqRVEKTUE0R0ExVUVDQk1IUW1WcGFtbHVaekVRTUE0R0ExVUVCeE1IUW1WcGFtbHVaekVmTUIwR0ExVUVDaE1XZDNndApiM0puTXk1amFHRnBibTFoYTJWeUxtOXlaekVTTUJBR0ExVUVDeE1KY205dmRDMWpaWEowTVNJd0lBWURWUVFECkV4bGpZUzUzZUMxdmNtY3pMbU5vWVdsdWJXRnJaWEl1YjNKbk1CNFhEVEl6TVRJd01UQTRORE14TkZvWERUSTQKTVRFeU9UQTRORE14TkZvd2dZOHhDekFKQmdOVkJBWVRBa05PTVJBd0RnWURWUVFJRXdkQ1pXbHFhVzVuTVJBdwpEZ1lEVlFRSEV3ZENaV2xxYVc1bk1SOHdIUVlEVlFRS0V4WjNlQzF2Y21jekxtTm9ZV2x1YldGclpYSXViM0puCk1RNHdEQVlEVlFRTEV3VmhaRzFwYmpFck1Da0dBMVVFQXhNaVlXUnRhVzR4TG5OcFoyNHVkM2d0YjNKbk15NWoKYUdGcGJtMWhhMlZ5TG05eVp6QlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJBM2pEZCs1d0N5Sgp2WW05Vko3eWJMU0R6STVzS3p2YXZOYkpvNVZJNU5GaERwZlpyTm1qN29hWlgzOTdjUkgzWWFvZnhlSnZ6ZnhKCnRHRDQxeEk3Vk9XamFqQm9NQTRHQTFVZER3RUIvd1FFQXdJR3dEQXBCZ05WSFE0RUlnUWcycklJNExwTDRmS1kKVGVORURrcWpMZWJ4M0k0U1lRMWtBbTJYM0JmNmpRSXdLd1lEVlIwakJDUXdJb0FneHFPVUhDOEdDTGxCZE44QgpHdE42S2RMVHJYZ09HSWJMMUxtUktMNlNWeGt3Q2dZSUtvWkl6ajBFQXdJRFNBQXdSUUloQVB4ZUR3Yis5Z2ZwCmdMdWgyWUJrV0NqbGVNdmg4NFNVUlcvN290QXFlOTVsQWlBeDMxdVUzbFZiUXJSWlhISmZldng3cVZ0am1ja3YKWmFDc1pnUWJhbmZueHc9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="

func TestBase64Encode(t *testing.T) {
	data := []byte("test")
	encoded := Base64Encode(data)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Errorf("Base64Decode failed: %v", err)
	}
	if string(decoded) != string(data) {
		t.Errorf("Base64Encode/Decode failed: expected %s, got %s", string(data), string(decoded))
	}
}

func TestBase64Decode(t *testing.T) {
	data := []byte("test")
	encoded := Base64Encode(data)
	decoded := Base64Decode(encoded)
	if string(decoded) != string(data) {
		t.Errorf("Base64Encode/Decode failed: expected %s, got %s", string(data), string(decoded))
	}
}

func TestCurrentMillSeconds(t *testing.T) {
	milliseconds := CurrentMillSeconds()
	if milliseconds <= 0 {
		t.Errorf("CurrentMillSeconds failed: expected positive number, got %d", milliseconds)
	}
}

func TestCurrentSeconds(t *testing.T) {
	seconds := CurrentSeconds()
	if seconds <= 0 {
		t.Errorf("CurrentSeconds failed: expected positive number, got %d", seconds)
	}
}

func TestPathExists(t *testing.T) {
	// Test case 1: Path exists
	path := "test1.txt"
	err := ioutil.WriteFile(path, []byte("test1"), 0600)
	if err != nil {
		t.Errorf("WriteFile failed: %v", err)
	}
	exists, err := PathExists(path)
	if err != nil {
		t.Errorf("PathExists failed: %v", err)
	}
	if !exists {
		t.Errorf("PathExists failed: expected true, got false")
	}
	err = os.Remove(path)
	if err != nil {
		t.Errorf("Remove failed: %v", err)
	}

	// Test case 2: Path does not exist
	path = "nonexistent.txt"
	exists, err = PathExists(path)
	if err != nil {
		t.Errorf("PathExists failed: %v", err)
	}
	if exists {
		t.Errorf("PathExists failed: expected false, got true")
	}
}

func TestParseCertificate(t *testing.T) {
	certBytes := []byte(TestCrt)
	cert, err := ParseCertificate(certBytes)
	if err != nil {
		t.Errorf("ParseCertificate failed: %v", err)
	}
	if cert == nil {
		t.Errorf("ParseCertificate failed: expected non-nil certificate, got nil")
	}
}

func TestCopy(t *testing.T) {
	from := map[string]interface{}{"key": "value"}
	to := map[string]interface{}{}
	err := Copy(&to, from)
	if err != nil {
		t.Errorf("Copy failed: %v", err)
	}
	if to["key"] != "value" {
		t.Errorf("Copy failed: expected value 'value', got %v", to["key"])
	}
}

func TestIsIgnoredFile(t *testing.T) {
	// Test case 1: File name contains ignored string
	fileName := "test.DS_Store"
	if !isIgnoredFile(fileName) {
		t.Errorf("isIgnoredFile failed: expected true, got false")
	}

	// Test case 2: File name does not contain ignored string
	fileName = "test.txt"
	if isIgnoredFile(fileName) {
		t.Errorf("isIgnoredFile failed: expected false, got true")
	}
}

func TestGetConfigShow(t *testing.T) {
	// Test case 1: GlobalConfig and ChainConf are nil
	config.GlobalConfig = nil
	if GetConfigShow() {
		t.Errorf("GetConfigShow failed: expected false, got true")
	}
}

func TestGetIsMainChain(t *testing.T) {
	// Test case 1: GlobalConfig and ChainConf are nil
	config.GlobalConfig = nil
	if GetIsMainChain() {
		t.Errorf("GetIsMainChain failed: expected false, got true")
	}
}

func TestCheckIsDIDContact(t *testing.T) {
	// Test case 1: Contract type is CMDID
	contractType := standard.ContractStandardNameCMDID
	if !CheckIsDIDContact(contractType) {
		t.Errorf("CheckIsDIDContact failed: expected true, got false")
	}

	// Test case 2: Contract type is not CMDID
	contractType = "other"
	if CheckIsDIDContact(contractType) {
		t.Errorf("CheckIsDIDContact failed: expected false, got true")
	}
}

// 测试SplicePath函数
func TestSplicePath(t *testing.T) {
	url := SplicePath("http://127.0.0.1:8888", "chain", "/transaction/", "123")
	if url != "http://127.0.0.1:8888/chain/transaction/123/" {
		t.Errorf("SplicePath failed: expected http://127.0.0.1:8888/chain/transaction/123, got %s", url)
	}
}

func TestCalculateSHA256(t *testing.T) {
	fileData := []byte("test data`")
	str := CalculateSHA256(fileData)
	if str != "a4dff0842b1663e64b74b1dea043033a032698a8b7491370f818b79d3ce477b9" {
		t.Errorf("CalculateSHA256 failed: expected a4dff0842b1663e64b74b1dea043033a032698a8b7491370f818b79d3ce477b9, got %s", str)
	}
}

func TestSplicePath1(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "basic paths",
			input:    []string{"http://example.com", "path", "to/resource"},
			expected: "http://example.com/path/to/resource/",
		},
		{
			name:     "with trailing slashes",
			input:    []string{"https://test.com/", "dir/", "file.txt"},
			expected: "https://test.com/dir/file.txt/",
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplicePath(tt.input...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBase64EncodeDecode(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"text data", []byte("Hello, World!")},
		{"binary data", []byte{0x00, 0x01, 0x02, 0x03, 0x7F, 0xFF}},
		{"empty data", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded := Base64Encode(tt.data)
			decoded := Base64Decode(encoded)

			assert.Equal(t, tt.data, decoded, "data should be equal after encode/decode")

			// 测试无效Base64
			invalidDecoded := Base64Decode("invalid!base64@string")
			assert.Nil(t, invalidDecoded, "invalid base64 should return nil")
		})
	}
}

func TestTimeFunctions(t *testing.T) {
	now := time.Now()

	t.Run("CurrentMillSeconds", func(t *testing.T) {
		millis := CurrentMillSeconds()
		assert.True(t, millis >= now.UnixNano()/1e6)
	})

	t.Run("CurrentSeconds", func(t *testing.T) {
		seconds := CurrentSeconds()
		assert.True(t, seconds >= now.Unix())
	})
}

func TestPathExists1(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	tests := []struct {
		name          string
		createPath    bool
		path          string
		expectedExist bool
	}{
		{"existing directory", true, tmpDir, true},
		{"non-existing", false, "/nonexistent/path/123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createPath {
				// 在临时目录中创建文件
				filePath := filepath.Join(tmpDir, "testfile.txt")
				f, err := os.Create(filePath)
				require.NoError(t, err)
				f.Close()

				exists, err := PathExists(filePath)
				assert.NoError(t, err)
				assert.True(t, exists)
			}

			exists, err := PathExists(tt.path)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedExist, exists)
		})
	}
}

func TestCertificateFunctions(t *testing.T) {
	// 生成测试证书
	certPEM, err := generateTestCertificate()
	require.NoError(t, err)

	t.Run("ParseCertificate", func(t *testing.T) {
		cert, err := ParseCertificate(certPEM)
		require.NoError(t, err)
		assert.Equal(t, "Test Certificate", cert.Subject.CommonName)

		// 测试无效证书
		_, err = ParseCertificate([]byte("invalid cert data"))
		assert.Error(t, err)
	})

	t.Run("X509CertToChainMakerCert", func(t *testing.T) {
		cert, err := ParseCertificate(certPEM)
		require.NoError(t, err)

		chainMakerCert, err := X509CertToChainMakerCert(cert)
		require.NoError(t, err)
		assert.Equal(t, cert.Subject.CommonName, chainMakerCert.Subject.CommonName)

		// 测试nil输入
		_, err = X509CertToChainMakerCert(nil)
		assert.Error(t, err)
	})
}

// 生成测试证书
func generateTestCertificate() ([]byte, error) {
	// 创建私钥
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Test Certificate",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	// 创建证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, err
	}

	// PEM编码
	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	return pemData, nil
}

func TestGenerateTableName(t *testing.T) {
	inputs := []ABIParamIn{
		{Name: "param1", Type: "uint256"},
		{Name: "param2", Type: "address"},
	}

	contractAddr := "0x1234567890abcdef"

	t.Run("with inputs", func(t *testing.T) {
		tableName := GenerateTableName(contractAddr, inputs)

		assert.True(t, strings.HasPrefix(tableName, "topic_0x123456"))
		assert.Len(t, tableName, len("topic_")+8+1+8) // topic_ + 8 chars + _ + 8 chars
	})
}

func TestGetDateFromTimestamp(t *testing.T) {
	t.Run("valid timestamp", func(t *testing.T) {
		timestamp := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC).Unix()
		result := GetDateFromTimestamp(timestamp)

		// 转换为东八区时间
		expected := time.Date(2023, 6, 15, 20, 0, 0, 0, time.FixedZone("CST", 8*60*60))

		assert.True(t, expected.Equal(result))
	})

	t.Run("invalid timestamp", func(t *testing.T) {
		result := GetDateFromTimestamp(-1)
		assert.True(t, result.IsZero())
	})
}

func TestGetAccountHashStr(t *testing.T) {
	_ = config.InitConfig("", "")
	expectedHash := "7c6a61c68ef8b9b6b061b28c348bc1ed7921cb53"
	accountHash := GetAccountHashStr()

	assert.Equal(t, expectedHash, accountHash, "Account hash should match the expected value")
}

func TestExtractAndHash7zFile(t *testing.T) {
	_ = config.InitConfig("", "")
	// 测试提取和哈希7z文件
	compressedData := []byte("7z archive data") // 假设这是压缩文件的数据
	_, _ = ExtractAndHash7zFile(compressedData)
}

func TestGetMd5Hash(t *testing.T) {
	_ = config.InitConfig("", "")
	_ = GetMd5Hash(1234567)
}

func TestParseContractABI(t *testing.T) {
	fileHeader := &multipart.FileHeader{}
	_, _ = ParseContractABI(fileHeader)
}
