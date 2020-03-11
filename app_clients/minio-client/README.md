# Minio-client for  application liveness an loadgen.
main.go is the primary execution script which runs in two mode based on environment variable passed.
- **liveness** mode which creates minio bucket and do the following operations
    - Loads data in minio bucket.
    - Unload minio bucket.
  For each successful load and unload operation client says `Liveness Running`.
- **loadgen** mode which simply writes amount of data provided as env-variable. `Note: The execution time of loadgen mode may vary depending on the amount of data to be written `

# Containerizing minio-client application.
- Keep main.go and docker file at **go-path** inside a same directory (ex: minio-client) and install the minio-client dependency by `GO111MODULE=on go get github.com/minio/minio-go/v6`.
- Then execute `CGO_ENABLED=0 go build main.go`. This will generate the go binary.
- Then using docker file build your own image for minio-client.
- Use minio-job.yaml and replace the image name name with your own. And finally tune the environment variables in **mini-job.yml** based on use-case.
- Finally apply **minio-job.yaml** in the cluster to benchmark mini0-application.