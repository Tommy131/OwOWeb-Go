/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-06 02:54:05
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-06 00:32:22
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package ocsp

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ocsp"
)

func OCSPHandler(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	req, err := ocsp.ParseRequest(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid OCSP request"})
		return
	}

	nonce, err := ExtractNonce(body)
	if err != nil {
		log.Printf("Failed to extract nonce: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to extract nonce"})
		return
	}

	respBytes, err := CreateOCSPResponse(userCert, caCert, req, userKey, nonce)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create OCSP response"})
		return
	}

	c.Data(http.StatusOK, "application/ocsp-response", respBytes)
}

// OCSP请求处理函数
func OCSPRequestHandler(c *gin.Context) {
	// 获取URL路径参数中的OCSP请求数据
	ocspRequestBase64 := c.Param("ocspRequest")

	// 解码Base64数据
	ocspRequestData, err := base64.StdEncoding.DecodeString(ocspRequestBase64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Base64 encoding"})
		return
	}

	// 解析OCSP请求
	ocspRequest, err := ocsp.ParseRequest(ocspRequestData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OCSP request"})
		return
	}

	log.Infof("Received an OCSP request: (%d)[%x][%x]", ocspRequest.SerialNumber, ocspRequest.IssuerNameHash, ocspRequest.IssuerKeyHash)

	// 生成OCSP响应
	ocspResponse, err := CreateOCSPResponse(caCert, caCert, ocspRequest, caKey, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OCSP response"})
		return
	}

	// 返回OCSP响应给客户端
	c.Data(http.StatusOK, "application/ocsp-response", ocspResponse)
}
