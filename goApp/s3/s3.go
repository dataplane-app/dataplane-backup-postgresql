package s3

import (
	"dataplane-backup/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client

func SetupInstance() {
	sess, err := minio.New(config.GConf.S3.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(config.GConf.S3.AccessKey, config.GConf.S3.SecureKey, ""),
		Secure: true,
	})

	if err != nil {
		log.Println(err)
	}

	// log.Printf("%#v\n", sess)

	Client = sess
}
