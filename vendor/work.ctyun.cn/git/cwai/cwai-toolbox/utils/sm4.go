package utils

import (
	"encoding/base64"
	"strings"

	"github.com/tjfoc/gmsm/sm4"

)

// key需要是16位的字符串
func Sm4Encrypt(key string, data string) string {
	encrypted, err := sm4.Sm4Ecb([]byte(key), []byte(data), true)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(encrypted)
}

func Sm4Decrypt(key string, encrypted string) string {
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return ""
	}
	decrypted, err := sm4.Sm4Ecb([]byte(key), encryptedByte, false)
	if err != nil {
		return ""
	}
	return string(decrypted)
}

func GenerateSm4Key(tenantID string) string {
	if len(tenantID) >= 16 {
		return tenantID[:16]
	} else {
		return tenantID + strings.Repeat("X", 16-len(tenantID))
	}
}

// Sm4DecryptByTenantID 以租户ID作为key按照SM4规则解析
func Sm4DecryptByTenantID(tenantID, encrypted string) string {
	return Sm4Decrypt(GenerateSm4Key(tenantID), encrypted)
}

// Sm4EncryptByTenantID 以租户ID作为key按照SM4规则加密
func Sm4EncryptByTenantID(tenantID, data string) string {
	return Sm4Encrypt(GenerateSm4Key(tenantID), data)
}