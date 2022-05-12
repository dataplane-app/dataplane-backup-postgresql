package s3

import (
	"dataplane-backup/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client = SetupInstance()

func SetupInstance() *minio.Client {
	sess, _ := minio.New(config.GConf.S3.Url, &minio.Options{
		Region: config.GConf.S3.Region,
		Creds:  credentials.NewStaticV4(config.GConf.S3.AccessKey, config.GConf.S3.SecureKey, ""),
	})
	return sess
}
