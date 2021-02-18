package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

const (
	accessKey   = "T3yaScJxNcWwqPN573gG3Iycb-1P9Wi3pfbAw2zf"
	secretKey   = "U1MSBJm8LFqjaCJp-EPvpAw33hZZoDCioKoMb1YF"
	bucket_name = "cuanpian"
)

func QiniuToken() string {
	putPolicy := storage.PutPolicy{
		Scope: bucket_name,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

func StorageSave(data []byte) (string, error) {
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuabei
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// putExtra := storage.PutExtra{
	// 	Params: map[string]string{
	// 		"x:name": "github logo",
	// 	},
	// }
	hash := md5.Sum(data)
	md5str := fmt.Sprintf("%x", hash) //将[]byte转成16进制
	key := bucket_name + "/" + md5str
	dataLen := int64(len(data))
	token := QiniuToken()
	fmt.Println(key, token)
	err := formUploader.Put(context.Background(), &ret, token, key, bytes.NewReader(data), dataLen, nil)
	if err != nil {
		return "", err
	}
	fmt.Println(ret.Key, ret.Hash)
	return ret.Key, nil
}
