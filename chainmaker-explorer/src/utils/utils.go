/*
Package utils comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"archive/zip"
	"bytes"
	"chainmaker_web/src/config"
	"crypto/md5"

	// #nosec G505
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/contract-utils/standard"
)

type ZipFileContract struct {
	SourcePath string
	SourceCode []byte
}

// SplicePath  函数用于拼接多个字符串，并返回拼接后的字符串
// eg: SplicePath("http://192.168.1.108:8080/", "chain1", "/transaction/", "6788795865a6d89e878ad9e999")
// get: http://192.168.1.108:8080/chain1/transaction/6788795865a6d89e878ad9e999/
func SplicePath(urls ...string) string {
	// 遍历传入的字符串数组 urls
	var bt bytes.Buffer
	// 将 url 写入 bt 中
	for _, url := range urls {
		// 如果 url 不以 / 结尾，则在其后添加 /
		bt.WriteString(url)
		if !strings.HasSuffix(url, "/") {
			bt.WriteString("/")
		}
		// 将 bt 中的字符串中的 // 替换为 /
	}
	// 将 bt 中的字符串中的 http:/ 替换为 http://
	r := strings.ReplaceAll(bt.String(), "//", "/")
	// 将 bt 中的字符串中的 https:/ 替换为 https://
	r = strings.ReplaceAll(r, "http:/", "http://")
	// 返回拼接后的字符串
	r = strings.ReplaceAll(r, "https:/", "https://")
	return r
}

// Base64Encode base64Encode
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 函数用于将Base64编码的字符串解码为字节数组
func Base64Decode(data string) []byte {
	// 使用base64标准库中的DecodeString函数将Base64编码的字符串解码为字节数组
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	// 如果解码过程中出现错误，则返回nil
	if err != nil {
		return nil
	}
	// 返回解码后的字节数组
	return decodeBytes
}

// CurrentMillSeconds cms
func CurrentMillSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// CurrentSeconds cs
func CurrentSeconds() int64 {
	return time.Now().UnixNano() / 1e9
}

// PathExists 函数用于判断指定路径是否存在
func PathExists(path string) (bool, error) {
	// 使用os.Stat函数获取指定路径的信息
	_, err := os.Stat(path)
	// 如果没有错误，说明路径存在
	if err == nil {
		return true, nil
	}
	// 如果错误类型为os.IsNotExist，说明路径不存在
	if os.IsNotExist(err) {
		return false, nil
	}
	// 否则返回错误
	return false, err
}

// ParseCertificate 解析证书
// @param certBytes 证书字节数组
// @return *x509.Certificate 证书
// @return error 错误
func ParseCertificate(certBytes []byte) (*x509.Certificate, error) {
	// 定义证书和错误变量
	var (
		cert *bcx509.Certificate
		err  error
	)
	// 解码证书
	block, rest := pem.Decode(certBytes)
	// 如果解码失败，则解析剩余的证书
	if block == nil {
		cert, err = bcx509.ParseCertificate(rest)
		// 否则解析解码后的证书
	} else {
		cert, err = bcx509.ParseCertificate(block.Bytes)
	}
	// 如果解析失败，则返回错误
	if err != nil {
		return nil, fmt.Errorf("[Parse cert] parseCertificate cert failed, %s", err)
	}

	// 返回解析后的证书
	return bcx509.ChainMakerCertToX509Cert(cert)
}

// X509CertToChainMakerCert x509 to cert
func X509CertToChainMakerCert(cert *x509.Certificate) (*bcx509.Certificate, error) {
	if cert == nil {
		return nil, fmt.Errorf("cert is nil")
	}
	der, err := bcx509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("fail to parse re-encode (marshal) public key in certificate: %v", err)
	}
	pk, err := asym.PublicKeyFromDER(der)
	if err != nil {
		return nil, fmt.Errorf("fail to parse re-encode (unmarshal) public key in certificate: %v", err)
	}
	newCert := &bcx509.Certificate{
		Raw:                         cert.Raw,
		RawTBSCertificate:           cert.RawTBSCertificate,
		RawSubjectPublicKeyInfo:     cert.RawSubjectPublicKeyInfo,
		RawSubject:                  cert.RawSubject,
		RawIssuer:                   cert.RawIssuer,
		Signature:                   cert.Signature,
		SignatureAlgorithm:          bcx509.SignatureAlgorithm(cert.SignatureAlgorithm),
		PublicKeyAlgorithm:          bcx509.PublicKeyAlgorithm(cert.PublicKeyAlgorithm),
		PublicKey:                   pk,
		Version:                     cert.Version,
		SerialNumber:                cert.SerialNumber,
		Issuer:                      cert.Issuer,
		Subject:                     cert.Subject,
		NotBefore:                   cert.NotBefore,
		NotAfter:                    cert.NotAfter,
		KeyUsage:                    cert.KeyUsage,
		Extensions:                  cert.Extensions,
		ExtraExtensions:             cert.ExtraExtensions,
		UnhandledCriticalExtensions: cert.UnhandledCriticalExtensions,
		ExtKeyUsage:                 cert.ExtKeyUsage,
		UnknownExtKeyUsage:          cert.UnknownExtKeyUsage,
		BasicConstraintsValid:       cert.BasicConstraintsValid,
		IsCA:                        cert.IsCA,
		MaxPathLen:                  cert.MaxPathLen,
		MaxPathLenZero:              cert.MaxPathLenZero,
		SubjectKeyId:                cert.SubjectKeyId,
		AuthorityKeyId:              cert.AuthorityKeyId,
		OCSPServer:                  cert.OCSPServer,
		IssuingCertificateURL:       cert.IssuingCertificateURL,
		DNSNames:                    cert.DNSNames,
		EmailAddresses:              cert.EmailAddresses,
		IPAddresses:                 cert.IPAddresses,
		URIs:                        cert.URIs,
		PermittedDNSDomainsCritical: cert.PermittedDNSDomainsCritical,
		PermittedDNSDomains:         cert.PermittedDNSDomains,
		ExcludedDNSDomains:          cert.ExcludedDNSDomains,
		PermittedIPRanges:           cert.PermittedIPRanges,
		ExcludedIPRanges:            cert.ExcludedIPRanges,
		PermittedEmailAddresses:     cert.PermittedEmailAddresses,
		ExcludedEmailAddresses:      cert.ExcludedEmailAddresses,
		PermittedURIDomains:         cert.PermittedURIDomains,
		ExcludedURIDomains:          cert.ExcludedURIDomains,
		CRLDistributionPoints:       cert.CRLDistributionPoints,
		PolicyIdentifiers:           cert.PolicyIdentifiers,
	}
	return newCert, nil
}

// Copy interface copy
func Copy(to, from interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("marshal from data err, %s", err.Error())
	}

	err = json.Unmarshal(b, to)
	if err != nil {
		return fmt.Errorf("unmarshal to data err, %s", err.Error())
	}

	return nil
}

// UnzipFileContractSource unzip file
// @param fileHeader *multipart.FileHeader file
// @return []*ZipFileContract  file list
// @return error error
func UnzipFileContractSource(fileHeader *multipart.FileHeader) ([]*ZipFileContract, error) {
	sourceFileList := make([]*ZipFileContract, 0)
	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return sourceFileList, err
	}
	defer file.Close()

	// 读取文件内容
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return sourceFileList, err
	}

	// 创建一个新的zip reader
	zipReader, err := zip.NewReader(bytes.NewReader(fileContent), int64(len(fileContent)))
	if err != nil {
		return sourceFileList, err
	}

	// 遍历压缩包中的每个文件
	for _, zipFile := range zipReader.File {
		// 过滤掉工具增加的文件
		if isIgnoredFile(zipFile.Name) {
			log.Infof("Skipping ignored file: %s\n", zipFile.Name)
			continue
		}

		// 打开压缩包中的文件
		zipFileReader, errOpen := zipFile.Open()
		if errOpen != nil {
			zipFileReader.Close()
			return sourceFileList, errOpen
		}

		// 读取压缩包中的文件内容
		zipFileContent, errRead := io.ReadAll(zipFileReader)
		if errRead != nil {
			zipFileReader.Close()
			return sourceFileList, errRead
		}

		// 检查文件内容是否为空
		if len(zipFileContent) == 0 {
			log.Infof("Skipping empty file: %s\n", zipFile.Name)
			// 关闭压缩包中的文件
			zipFileReader.Close()
			continue
		}

		// 创建 ContractSourceFile 实例
		sourceFile := &ZipFileContract{
			SourcePath: zipFile.Name, // 使用压缩包内部的文件路径
			SourceCode: zipFileContent,
		}

		sourceFileList = append(sourceFileList, sourceFile)

		// 关闭压缩包中的文件
		zipFileReader.Close()
	}

	return sourceFileList, nil
}

// isIgnoredFile 检查文件名是否是被忽略的文件
// @param fileName string 文件名
// @return bool 是否被忽略
func isIgnoredFile(fileName string) bool {
	// 定义被忽略的文件名列表
	ignoredFiles := []string{
		".DS_Store",
		"__MACOSX",
		"Thumbs.db",
		"desktop.ini",
	}

	// 遍历被忽略的文件名列表
	for _, ignored := range ignoredFiles {
		// 如果文件名包含被忽略的文件名，则返回true
		if strings.Contains(fileName, ignored) {
			return true
		}
	}
	// 如果文件名不包含被忽略的文件名，则返回false
	return false
}

// GetConfigShow - 获取链配置是否显示
func GetConfigShow() bool {
	if config.GlobalConfig == nil || config.GlobalConfig.ChainConf == nil {
		return false
	}

	return config.GlobalConfig.ChainConf.ShowConfig
}

// GetIsMainChain - 是否是主链
func GetIsMainChain() bool {
	if config.GlobalConfig == nil || config.GlobalConfig.ChainConf == nil {
		return false
	}

	return config.GlobalConfig.ChainConf.IsMainChain
}

// IsDIDContact 判断合约类型是否为CMDID
func CheckIsDIDContact(contractType string) bool {
	// 如果合约类型等于标准合约类型CMDID，则返回true
	return contractType == standard.ContractStandardNameCMDID
}

// ExtractAndHash7zFile 解压 7z 文件并计算文件的 sha256 哈希,注意：此函数依赖于 7z 命令行工具
// @param compressedData 压缩的 7z 文件数据
// @return 文件的 sha256 哈希值
// @return 错误信息
func ExtractAndHash7zFile(compressedData []byte) (string, error) {
	// 使用 os.CreateTemp 创建临时文件存储压缩数据
	tmpFile, err := os.CreateTemp("", "temp-file-*.7z")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name()) // 确保临时文件在函数结束时被删除

	// 将压缩数据写入临时文件
	_, err = tmpFile.Write(compressedData)
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	// 使用 7z 解压该文件
	cmd := exec.Command("7z", "x", tmpFile.Name(), "-so") // -so: 从标准输入读取压缩数据并解压
	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行解压命令
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to extract 7z data: %v", err)
	}

	// 计算并返回该文件的 sha256 哈希值,
	hashValue := CalculateSHA256(out.Bytes())
	log.Infof("ExtractAndHash7zFile ====out.Bytes hashValue = %v", hashValue)
	return hashValue, nil
}

// CalculateSHA256 计算文件的 sha256 哈希值
// @param fileData 文件数据
// @return 文件的 sha256 哈希值
func CalculateSHA256(fileData []byte) string {
	// 计算 sha256 哈希
	hash := sha256.New()
	hash.Write(fileData)
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed) // 将哈希值转换为十六进制字符串并返回
}

// GetMd5Hash 获取MD5哈希值
// @param randomNum 随机数
// @return string MD5哈希值
func GetMd5Hash(randomNum int64) string {
	//获取账户密码
	webConf := config.GlobalConfig.WebConf
	password := webConf.AdminPassword
	randomStr := strconv.FormatInt(randomNum, 10)
	passwordRandomNum := password + "_" + randomStr
	hasher := md5.New()
	hasher.Write([]byte(passwordRandomNum))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetAccountHashStr 获取账户密码的哈希值,用于登录验证
// @return string 哈希值
func GetAccountHashStr() string {
	//获取账户密码
	webConf := config.GlobalConfig.WebConf
	password := webConf.AdminPassword

	// 创建一个新的hash.Hash256哈希实例
	//hash := sha256.New()
	hash := sha1.New()

	// 写入要哈希的数据
	hash.Write([]byte(password))

	// 计算哈希值
	hashBytes := hash.Sum(nil)

	// 将字节转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

// ParseContractABI 解析合约ABI
func ParseContractABI(fileHeader *multipart.FileHeader) ([]*ContractABI, error) {
	abiList := make([]*ContractABI, 0)

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		return abiList, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return abiList, fmt.Errorf("failed to read file: %v", err)
	}

	// 解析JSON内容到ContractABI结构体切片
	err = json.Unmarshal(fileContent, &abiList)
	if err != nil {
		return abiList, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return abiList, nil
}

// 生成表名（新格式）
func GenerateTableName(contractAddr string, inputs []ABIParamIn) string {
	// 生成输入参数哈希
	inputHash, err := generateInputHash(inputs)
	if err != nil {
		log.Errorf("generateInputHash err: %s", err)
		return ""
	}

	addrPart := contractAddr
	if len(addrPart) > 8 {
		addrPart = addrPart[:8]
	}
	hashPart := inputHash
	if len(inputHash) > 8 {
		hashPart = inputHash[:8]
	}
	return fmt.Sprintf("topic_%s_%s", addrPart, hashPart)
}

// 生成输入参数哈希
func generateInputHash(inputs []ABIParamIn) (string, error) {
	jsonData, err := json.Marshal(inputs)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(jsonData)), nil
}

// ConvertTimestampToGMT8 将秒级时间戳转换为东八区日期
func GetDateFromTimestamp(timestamp int64) time.Time {
	// 处理无效时间戳，返回零值时间
	if timestamp <= 0 {
		return time.Time{}
	}

	// 1. 从UTC时间戳创建时间对象
	t := time.Unix(timestamp, 0) // UTC时间

	// 2. 转换为东八区时间
	// 方法1: 使用预定义的上海时区（推荐）
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 方法2: 如果时区数据不可用，创建固定时区
		loc = time.FixedZone("GMT+8", 8*60*60)
	}

	// 3. 转换为东八区时间
	return t.In(loc)
}
