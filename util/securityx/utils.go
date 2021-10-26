package securityx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"hash"
	"io"
	"strings"
)

func Md5(str string, salt ...string) string {
	var _salt string

	if len(salt) > 0 && salt[0] != "" {
		_salt = salt[0]
	}

	h := md5.New()
	h.Write([]byte(str))
	s1 := hex.EncodeToString(h.Sum(nil))

	if _salt == "" {
		return s1
	}

	return Md5(s1 + _salt)
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) string {
	buf, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		return ""
	}

	return string(buf)
}

func Sha1(str string) (string, error) {
	return shaStr(sha1.New(), []byte(str))
}

func Sha256(str string) (string, error) {
	return shaStr(sha256.New(), []byte(str))
}

func Sha512(str string) (string, error) {
	return shaStr(sha512.New(), []byte(str))
}

func HmacSha1(str string, keyBuf []byte) (string, error) {
	return hmacShaString(sha1.New, str, keyBuf)
}

func HmacSha256(str string, keyBuf []byte) (string, error) {
	return hmacShaString(sha256.New, str, keyBuf)
}

func HmacSha512(str string, keyBuf []byte) (string, error) {
	return hmacShaString(sha512.New, str, keyBuf)
}

func RsaEncrypt(str string, publicKeyBuf []byte) (ret string, err error) {
	block, _ := pem.Decode(publicKeyBuf)

	if block == nil {
		err = errors.New("invalid rsa public key")
		return
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return
	}

	publicKey, ok := pub.(*rsa.PublicKey)

	if !ok {
		err = errors.New("invalid rsa public key")
		return
	}

	buf, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(str))

	if err != nil {
		return
	}

	ret = string(buf)
	return
}

func RsaDecrypt(str string, privateKeyBuf []byte) (ret string, err error) {
	block, _ := pem.Decode(privateKeyBuf)

	if block == nil {
		err = errors.New("invalid rsa private key")
		return
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return
	}

	buf, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(str))

	if err != nil {
		return
	}

	ret = string(buf)
	return
}

func AesCbcEncrypt(str string, keyBuf []byte, paddingMode ...string) (ret string, err error) {
	block, err := aes.NewCipher(keyBuf)

	if err != nil {
		return
	}

	blockSize := block.BlockSize()
	_paddingMode := "pkcs7"

	if len(paddingMode) > 0 && paddingMode[0] != "" {
		_paddingMode = strings.ToLower(paddingMode[0])
	}

	srcBuf := []byte(str)

	switch _paddingMode {
	case "pkcs7":
		srcBuf = pkcs7Padding(srcBuf, blockSize)
	case "pkcs5":
		srcBuf = pkcs5Padding(srcBuf, blockSize)
	default:
		err = errors.New("unsupported padding mode")
		return
	}

	dstBuf := make([]byte, blockSize + len(srcBuf))
	ivBuf := dstBuf[:blockSize]
	_, err = io.ReadFull(rand.Reader, ivBuf)

	if err != nil {
		return
	}

	mode := cipher.NewCBCDecrypter(block, ivBuf)
	mode.CryptBlocks(dstBuf[blockSize:], srcBuf)
	ret = Base64Encode(string(dstBuf))
	return
}

func AesCbcDecrypt(str string, keyBuf, ivBuf []byte) (ret string, err error) {
	block, err := aes.NewCipher(keyBuf)

	if err != nil {
		return
	}

	srcBuf := []byte(str)
	dstBuf := make([]byte, len(srcBuf))
	mode := cipher.NewCBCDecrypter(block, ivBuf)
	mode.CryptBlocks(dstBuf, srcBuf)
	ret = string(dstBuf)
	return
}

func hmacShaString(fn func() hash.Hash, str string, keyBuf []byte) (ret string, err error) {
	h := hmac.New(fn, keyBuf)
	_, err = h.Write([]byte(str))

	if err != nil {
		return
	}

	ret = hex.EncodeToString(h.Sum(nil))
	return
}

func shaStr(h hash.Hash, buf []byte) (ret string, err error) {
	_, err = h.Write(buf)

	if err != nil {
		return
	}

	ret = hex.EncodeToString(h.Sum([]byte(nil)))
	return
}

func pkcs7Padding(buf []byte, blockSize int) []byte {
	n1 := blockSize - len(buf) % blockSize
	padBuf := bytes.Repeat([]byte{byte(n1)}, n1)
	return append(buf, padBuf...)
}

func pkcs7Unpadding(buf []byte) []byte {
	n1 := len(buf)
	n2 := int(buf[n1 - 1])
	return buf[:(n1 - n2)]
}

func pkcs5Padding(buf []byte, blockSize int) []byte {
	n1 := blockSize - len(buf) % blockSize
	padBuf := bytes.Repeat([]byte{byte(n1)}, n1)
	return append(buf, padBuf...)
}

func pkcs5Unpadding(buf []byte) []byte {
	n1 := len(buf)
	n2 := int(buf[n1 - 1])
	return buf[:(n1 - n2)]
}
