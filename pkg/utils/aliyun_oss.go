package utils

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"strings"
)

// Upload2AliyunOss creates a new object in aliyun OSS and it will overwrite the original one if it exists already.
//
// objectKey    the object key in UTF-8 encoding. The length must be between 1 and 1023, and cannot start with "/" or "\".
// reader    io.Reader instance for reading the data for uploading
func UploadQrcode2AliyunOss(objectKey string, reader io.Reader) (string, error) {

	var bucketName = Config.GetString("storage.aliyun.bucketname4qrcode")
	if bucketName == "" {
		Log.Errorln("Wrong configuration: storage.aliyun.bucketname4qrcode is empty")
		return "", errors.New("storage.aliyun.bucketname4qrcode is empty")
	}
	options := []oss.Option{
		oss.ObjectACL(oss.ACLPublicRead),
	}
	return upload2AliyunOss(objectKey, bucketName, reader, options)
}

func upload2AliyunOss(objectKey string, bucketName string, reader io.Reader, options []oss.Option) (string, error) {
	var endpoint = Config.GetString("storage.aliyun.endpoint")
	if endpoint == "" {
		Log.Errorln("Wrong configuration: storage.aliyun.endpoint is empty")
		return "", errors.New("storage.aliyun.endpoint is empty")
	}
	var accessKeyId = Config.GetString("storage.aliyun.accesskeyid")
	if accessKeyId == "" {
		Log.Errorln("Wrong configuration: storage.aliyun.accesskeyid is empty")
		return "", errors.New("storage.aliyun.accesskeyid is empty")
	}
	var accessKeySecret = Config.GetString("storage.aliyun.accesskeysecret")
	if accessKeySecret == "" {
		Log.Errorln("Wrong configuration: storage.aliyun.accesskeysecret is empty")
		return "", errors.New("storage.aliyun.accesskeysecret is empty")
	}

	// 创建OSSClient实例
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return "", err
	}

	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}

	// 上传文件
	err = bucket.PutObject(objectKey, reader, options...)
	if err != nil {
		return "", err
	}

	// endpoint example: http://oss-ap-southeast-1.aliyuncs.com
	// objUrl example: http://yuudidi-qrcode-test.oss-ap-southeast-1.aliyuncs.com/123.png
	var objUrl string
	if strings.HasPrefix(endpoint, "http://") {
		objUrl = "http://" + bucketName + "." + endpoint[len("http://"):] + "/" + objectKey
	} else if strings.HasPrefix(endpoint, "https://") {
		objUrl = "https://" + bucketName + "." + endpoint[len("https://"):] + "/" + objectKey
	} else {
		return "", errors.New("endpoint is invalid")
	}

	return objUrl, nil
}
