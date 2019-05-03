package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"log"

	"net/url"
	"os"
	//	"sort"
	//	"strings"
	//	"time"
)

func Sign(data []byte) (signature string, err error) {
	var h hash.Hash
	var hType crypto.Hash

	h = sha256.New()
	hType = crypto.SHA256
	h.Write(data)
	d := h.Sum(nil)

	privateKey, err := ioutil.ReadFile("alipay_private.pem")
	if err != nil {
		log.Fatalln(err)
	}
	pri, err := GetPriKey(privateKey)
	bs, err := rsa.SignPKCS1v15(rand.Reader, pri, hType, d)
	if err != nil {
		fmt.Println(err)
		return
	}
	signature = base64.StdEncoding.EncodeToString(bs)
	//fmt.Println(signature)
	return
}

func VerifySign(data []byte, sign string) (signature string, err error) {
	var h hash.Hash
	var hType crypto.Hash
	h = sha256.New()
	hType = crypto.SHA256
	h.Write(data)
	d := h.Sum(nil)
	key, err := ioutil.ReadFile("alipay_public.pem")
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	pub, err := GetPubKey(key)
	if err != nil {
		fmt.Println(err)
		return
	}
	s, _ := base64.StdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(pub, hType, d, s)
	if err != nil {
		fmt.Println("====》》》》》》", err)
	} else {
		fmt.Println("verify ok")
	}
	return
}

func GenRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	/* 核心代码开始 */
	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: derStream,
	}
	fi, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(fi, block)
	if err != nil {
		return err
	}
	/* 核心代码结束 */
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	/* 核心代码开始 */
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	fi, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(fi, block)
	/* 核心代码结束 */
	return err
}

//// 全局变量
var privateKey, publicKey []byte

/**
 * 功能：获取RSA公钥长度
 * 参数：public
 * 返回：成功则返回 RSA 公钥长度，失败返回 error 错误信息
 */
func GetPubKeyLen(pubKey []byte) (int, error) {
	if pubKey == nil {
		return 0, errors.New("input arguments error")
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		fmt.Println("========")
		return 0, errors.New("public rsaKey error")
	}
	fmt.Println("========", block.Bytes)
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return 0, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	fmt.Println(pub)
	return pub.N.BitLen(), nil
}

/**
 * 功能：获取RSA公钥长度
 * 参数：public
 * 返回：成功则返回 RSA 公钥长度，失败返回 error 错误信息
 */
func GetPubKey(pubKey []byte) (*rsa.PublicKey, error) {
	if pubKey == nil {
		return nil, errors.New("input arguments error")
	}
	block, _ := pem.Decode(pubKey)
	if block == nil {
		fmt.Println("========")
		return nil, errors.New("public rsaKey error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return pub, nil

}

/*
   获取RSA私钥长度
   PriKey
   成功返回 RSA 私钥长度，失败返回error
*/
func GetPriKeyLen(priKey []byte) (int, error) {
	if priKey == nil {
		return 0, errors.New("input arguments error")
	}
	block, _ := pem.Decode(priKey)
	if block == nil {
		return 0, errors.New("private rsaKey error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return 0, err
	}

	return priv.N.BitLen(), nil
}

/*
   获取RSA私钥长度
   PriKey
   成功返回 RSA 私钥长度，失败返回error
*/
func GetPriKey(priKey []byte) (*rsa.PrivateKey, error) {
	if priKey == nil {
		return nil, errors.New("input arguments error")
	}
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private rsaKey error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func TT() {
	gourl := "https://openapi.alipay.com/gateway.do?app_id=2019043064381424&charset=utf-8&code=62b78b4abc5140fcaa95ab08f19aZX52&grant_type=authorization_code&method=alipay.system.oauth.token&sign=mLJ2WmURQ%2F8gMMyS6T3J%2BxMzijt6TDEcmie%2FkxPmHyLjh%2BCDMl2Glm%2F%2BovYD0SVbRDhd8K%2FvDqhT%2FgQ0ZlRZEgmQJCyoq%2Bw79HjVmgM5OYZ0UXBeSfSkAThFiSlR4fuuybkD8i%2BHUGtsboFvKlKK0Tz3mQkQY0kso8nmgBwr1GB2XAOAXXlg42SJ87Mi6FY6G7hUaXG5cUcLsR9QsT2t9vyljf3D5HmzC8Zy8Gd4Nj%2BHZOhsEYFZ9nGNIoL0xcIpMcT%2FoMcmpEA2KOhl9JgaIlYbtBqVSMAAozZe3KOvt1DWUHzE9qdtpgPdr0QQZ3kQ0IcfdpfSBHjhY4Uz26cTfw%3D%3D&sign_type=RSA2&timestamp=2019-04-30+23%3A55%3A40&version=1.0"

	d, _ := url.QueryUnescape(gourl)
	fmt.Println("ssssssss====>", d)

	//	gourl = strings.Replace(gourl, "https://openapi.alipay.com/gateway.do?", "", -1)
	keyMap, _ := url.ParseQuery(gourl)
	data := url.Values{}

	signstr1, _ := url.QueryUnescape(keyMap.Encode())
	fmt.Println("--------", signstr1)

	sign := keyMap["sign"]

	for k, v := range keyMap {
		if k == "sign" || k == "sign_type" {
			continue
		}
		data.Set(k, v[0])
	}

	signstr, _ := url.QueryUnescape(data.Encode())

	fmt.Println("========", sign[0])
	fmt.Println("app_id=2019043064381424&charset=utf-8&code=62b78b4abc5140fcaa95ab08f19aZX52&grant_type=authorization_code&method=alipay.system.oauth.token&timestamp=2019-04-30 23:55:40&sign_type=RSA2&version=1.0", sign[0])

	VerifySign([]byte(signstr), sign[0])

}
