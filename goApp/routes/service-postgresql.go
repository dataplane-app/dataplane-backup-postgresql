package routes

import (
	"bytes"
	"dataplane-backup/config"
	"dataplane-backup/s3"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

type DatabaseBackup struct {
	Status      string `json:"status,omitempty"`
	S3Bucket    string `json:"s3_bucket,omitempty"`
	S3Size      int64  `json:"s3_size,omitempty"`
	S3ObjectKey string `json:"s3_object_key,omitempty"`
	S3Response  string
}

func RunPostgresBackup(c *fiber.Ctx) error {

	var err error

	ctx := c.Context()

	// Run the back up command::
	dumpFilename := fmt.Sprintf("/app/backup/%v-db-timescaledb-2.5.1-pg14.sql", time.Now().Unix())
	dbConfig := config.GConf.PostgresDatabase

	command := fmt.Sprintf(
		"PGPASSWORD=%s PGSSLMODE=%s pg_dump --create --clean -h %s -p %s -U %s > %s",
		dbConfig.Password, dbConfig.SSL, dbConfig.Host, dbConfig.Port, dbConfig.User, dumpFilename,
	)

	// PGPASSWORD=xxx pg_dump -h 127.0.0.1 --create --clean --port=5000 -U postgres dataplane > $(date +"%Y-%m-%d-%H")-dataplane-demo-db-timescaledb-2.5.1-pg14-.sql

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	if err != nil {
		// cleanup files if something goes wrong
		os.Remove(dumpFilename)
		fmt.Println(outb.String())
		fmt.Println(errb.String())
		return err
	}

	// Compress the file
	filename := fmt.Sprintf("%s.tar.gz", dumpFilename)

	cmdn := exec.CommandContext(ctx, "tar", "-czf", filename, dumpFilename)
	cmdn.Stdout = &outb
	cmdn.Stderr = &errb
	err = cmdn.Run()
	if err != nil {
		// cleanup files if something goes wrong
		os.Remove(filename)
		fmt.Println(outb.String())
		fmt.Println(errb.String())
		return err
	}

	info, err := s3.Client.FPutObject(
		ctx,
		config.GConf.S3.Bucket,
		fmt.Sprintf("backup/database/%s", filename),
		filename,
		minio.PutObjectOptions{
			ContentType: "application/x-tar",
		},
	)
	if err != nil {
		return err
	}

	// Clean up
	os.Remove(dumpFilename)
	os.Remove(filename)

	// --- send response
	log.Println("postgres backup complete", "size:", info.Size)

	b, _ := json.Marshal(info)

	var dbBackup DatabaseBackup
	dbBackup.S3Response = string(b)
	dbBackup.S3Bucket = info.Bucket
	dbBackup.S3ObjectKey = info.Key
	dbBackup.S3Size = info.Size

	return c.JSON(dbBackup)

}
