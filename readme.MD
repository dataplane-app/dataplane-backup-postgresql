### APIs -
`POST` to `/postgres-backup` starts a Postgres Database Backup to S3.


### Start the server
```shell
go run server.go
```

### Folders
* /app - application
* /app/backup - backup files

### Environment variables

Database
* DB_USER
* DB_DATABASE
* DB_HOST
* DB_PWD
* DB_PORT
* DB_SSL

S3 API
* S3_ACCESS_KEY
* S3_SECRET
* S3_BUCKET
* S3_REGION
* S3_URL

### Build docker images
```shell
export dpversion=0.0.x

docker build -t dataplane/dataplane-backup-postgresql-s3:$dpversion -f Dockerfile.alpine .
docker tag dataplane/dataplane-backup-postgresql-s3:$dpversion dataplane/dataplane-backup-postgresql-s3:$dpversion
docker push dataplane/dataplane-backup-postgresql-s3:$dpversion
docker tag dataplane/dataplane-backup-postgresql-s3:$dpversion dataplane/dataplane-backup-postgresql-s3:latest
docker push dataplane/dataplane-backup-postgresql-s3:latest
```
