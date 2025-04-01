/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-06 02:24:56
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2025-04-01 19:10:13
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package ocsp

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"
	"owoweb/utils"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ocsp"
	"golang.org/x/term"
)

// 使用一个结构体存储所有文件路径
type CertPaths struct {
	CACertFile   string
	CAKeyFile    string
	UserCertFile string
	UserKeyFile  string
	CRLFile      string
}

// 定义常量
const (
	Good         = 0
	Revoked      = 1
	Unknown      = 2
	ServerFailed = 3
)

// 映射从 int 到 string
var intToString = map[int]string{
	Good:         "Good",
	Revoked:      "Revoked",
	Unknown:      "Unknown",
	ServerFailed: "ServerFailed",
}

// 映射从 string 到 int
var stringToInt = map[string]int{
	"Good":         Good,
	"Revoked":      Revoked,
	"Unknown":      Unknown,
	"ServerFailed": ServerFailed,
}

// 从 int 转换为 string
func IntToStr(code int) (string, bool) {
	str, ok := intToString[code]
	return str, ok
}

// 从 string 转换为 int
func StrToInt(s string) (int, bool) {
	code, ok := stringToInt[s]
	return code, ok
}

var nonceOID = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 48, 1, 2}

var caCert *x509.Certificate
var caKey crypto.Signer
var userCert *x509.Certificate
var userKey crypto.Signer
var crl *x509.RevocationList

func init() {
	// 定义文件夹路径常量
	const ocspPath = utils.STORAGE_PATH + "ocsp/"
	var certPaths = CertPaths{
		CACertFile:   ocspPath + "OwOTeam_Root_CA.crt",
		CAKeyFile:    ocspPath + "OwOTeam_Root_CA.key",
		UserCertFile: ocspPath + "owoserver.com.crt",
		UserKeyFile:  ocspPath + "owoserver.com.key",
		CRLFile:      ocspPath + "rootca.crl",
	}

	var err error
	caCert, err = LoadCertificate(certPaths.CACertFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}

	caKey, err = LoadKey(certPaths.CAKeyFile)
	if err != nil {
		log.Fatalf("Failed to load CA key: %v", err)
	}

	userCert, err = LoadCertificate(certPaths.UserCertFile)
	if err != nil {
		log.Fatalf("Failed to load user certificate: %v", err)
	}

	userKey, err = LoadKey(certPaths.UserKeyFile)
	if err != nil {
		log.Fatalf("Failed to load user key: %v", err)
	}

	crl, err = LoadCRL(certPaths.CRLFile)
	if err != nil {
		log.Fatalf("Failed to load CRL: %v", err)
	}

	// 初始化数据库
	SetupDatabase(utils.DATABASE_PATH + "ocsp_database.db")
	// 检查证书吊销状态
	CheckCRL()

	log.Println("Loaded OCSP Services.")
}

// LoadCertificate loads a PEM-encoded certificate from a file
func LoadCertificate(filename string) (*x509.Certificate, error) {
	certPEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(certPEM)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// LoadKey loads a PEM-encoded private key from a file
func LoadKey(filename string) (crypto.Signer, error) {
	keyPEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing private key")
	}

	var key crypto.Signer
	switch block.Type {
	case "PRIVATE KEY":
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			// If parsing as PKCS1 fails, try parsing as PKCS8
			keyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return nil, err
			}
			key, ok := keyInterface.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("not an RSA private key")
			}
			return key, nil
		}
		key = privateKey

	case "RSA PRIVATE KEY":
		// 提示用户输入密码
		fmt.Print("Please enter CA Certificate's password: ")
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("Error reading password: %v", err)
		}
		fmt.Println()

		// 解密PEM块
		decryptedKeyBytes, err := x509.DecryptPEMBlock(block, password)
		if err != nil {
			log.Fatalf("Failed to decrypt PEM block: %v", err)
		}

		// 解析解密后的私钥
		privateKey, err := x509.ParsePKCS1PrivateKey(decryptedKeyBytes)
		if err != nil {
			log.Fatalf("Failed to parse RSA private key: %v", err)
		}
		key = privateKey

	case "EC PRIVATE KEY":
		ecKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		key = ecKey

	default:
		return nil, errors.New("unsupported key type")
	}

	return key, nil
}

// LoadCRL loads and parses a CRL file
func LoadCRL(filename string) (*x509.RevocationList, error) {
	crlPEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(crlPEM)
	if block == nil || block.Type != "X509 CRL" {
		return nil, errors.New("failed to decode PEM block containing CRL")
	}

	crl, err := x509.ParseRevocationList(block.Bytes)
	if err != nil {
		return nil, err
	}
	return crl, nil
}

// CheckCRL get information of Certificate Revoke List
func CheckCRL() {
	// 输出 CRL 签发者信息
	log.Infof("Issuer: %v", crl.Issuer)

	// 输出吊销的证书数量
	log.Infof("Amount of revoked Certificates: %d", len(crl.RevokedCertificateEntries))

	// 列出吊销的证书信息
	for _, revokedCert := range crl.RevokedCertificateEntries {
		log.Infof("Revoked Serial Number: %v, Revoked Time: %v", revokedCert.SerialNumber, revokedCert.RevocationTime)
	}
}

// IsCertificateRevoked checks if the certificate is revoked based on the CRL
func IsCertificateRevoked(serial *big.Int) bool {
	for _, revoked := range crl.RevokedCertificateEntries {
		if revoked.SerialNumber.Cmp(serial) == 0 {
			return true
		}
	}
	return false
}

// ExtractNonce extracts the nonce from the ASN.1 encoded OCSP request
func ExtractNonce(reqBytes []byte) ([]byte, error) {
	var ocspRequest ocspRequest
	_, err := asn1.Unmarshal(reqBytes, &ocspRequest)
	if err != nil {
		return nil, err
	}

	for _, ext := range ocspRequest.TBSRequest.RequestExtensions {
		if ext.Id.Equal(nonceOID) {
			return ext.Value, nil
		}
	}
	return nil, nil
}

// CreateOCSPResponse creates an OCSP response for the given request
func CreateOCSPResponse(cert *x509.Certificate, issuer *x509.Certificate, req *ocsp.Request, key crypto.Signer, nonce []byte) ([]byte, error) {
	// 获取证书状态
	certStatus, err := GetCertificateStatus(req.SerialNumber.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get certificate status: %v", err)
	}

	status, _ := StrToInt(certStatus.Status)

	template := ocsp.Response{
		Status:       status,
		SerialNumber: req.SerialNumber,
		ThisUpdate:   time.Now(),
		NextUpdate:   time.Now().Add(24 * time.Hour),
		Certificate:  cert,
	}

	// 如果证书被吊销，设置吊销时间和原因
	if status == ocsp.Revoked {
		template.RevokedAt = certStatus.RevokedAt
		template.RevocationReason = certStatus.RevocationReason
	}

	if cert == nil {
		return nil, errors.New("cert is nil")
	}
	if issuer == nil {
		return nil, errors.New("issuer is nil")
	}
	if key == nil {
		return nil, errors.New("key is nil")
	}

	// Include the nonce in the response if present in the request
	if nonce != nil {
		template.Extensions = []pkix.Extension{
			{
				Id:       nonceOID,
				Critical: false,
				Value:    nonce,
			},
		}
	}

	respBytes, err := ocsp.CreateResponse(issuer, cert, template, key)
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}
