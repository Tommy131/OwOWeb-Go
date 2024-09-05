/*
 *        _____   _          __  _____   _____   _       _____   _____
 *      /  _  \ | |        / / /  _  \ |  _  \ | |     /  _  \ /  ___|
 *      | | | | | |  __   / /  | | | | | |_| | | |     | | | | | |
 *      | | | | | | /  | / /   | | | | |  _  { | |     | | | | | |   _
 *      | |_| | | |/   |/ /    | |_| | | |_| | | |___  | |_| | | |_| |
 *      \_____/ |___/|___/     \_____/ |_____/ |_____| \_____/ \_____/
 *
 *  Copyright (c) 2023 by OwOTeam-DGMT (OwOBlog).
 * @Date         : 2024-06-30 23:01:04
 * @Author       : HanskiJay
 * @LastEditors  : HanskiJay
 * @LastEditTime : 2024-09-06 01:06:09
 * @E-Mail       : support@owoblog.com
 * @Telegram     : https://t.me/HanskiJay
 * @GitHub       : https://github.com/Tommy131
 */
package ocsp

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"time"
)

// CertificateStatus 表示证书的状态和相关信息
type CertificateStatus struct {
	SerialNumber     string `gorm:"primaryKey; not null"` // 将 serial_number 设为主键
	CommonName       string `gorm:"type:varchar(100); not null"`
	OrganizationUnit string `gorm:"type:varchar(100); default:General"`
	Organization     string `gorm:"type:varchar(100); default:OwOTeam"`
	State            string `gorm:"type:varchar(100); default:Bayern"`
	Country          string `gorm:"type:varchar(10); default:DE"`
	Email            string `gorm:"type:varchar(100); default: support@owoblog.com"`
	Status           string `gorm:"type:varchar(20); not null"`
	RevokedAt        time.Time
	RevocationReason int       `gorm:"default:0"`
	IssuedDate       time.Time `gorm:"autoCreateTime"`
	ExpiredDate      time.Time
}

// TableName 方法返回自定义的表名
func (CertificateStatus) TableName() string {
	return "certificate_status"
}

// 定义 CertID 结构体
type CertID struct {
	HashAlgorithm  pkix.AlgorithmIdentifier
	IssuerNameHash asn1.RawValue
	IssuerKeyHash  asn1.RawValue
	SerialNumber   asn1.RawValue
}

// 定义 RequestListEntry 结构体
type RequestListEntry struct {
	CertID CertID
}

// 定义 TBSRequest 结构体
type TBSRequest struct {
	Version           int `asn1:"optional,explicit,default:0,tag:0"`
	RequestList       []RequestListEntry
	RequestExtensions []pkix.Extension `asn1:"optional,explicit,tag:2"`
}

// 定义 ocspRequest 结构体
type ocspRequest struct {
	TBSRequest TBSRequest
}
