package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	minio "github.com/minio/minio-go"
)

var (
	accessKeyID           string
	secretAccessKey       string
	serviceName           string
	namespace             string
	port                  string
	livenessCheckInterval int
	livenessRetryCount    int
	mode                  string
	capacity              int
)

func main() {

	err := initVars()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileName := "bigfile.test"
	bucketName := "test-bucket" + "-" + fmt.Sprintf("%d", time.Now().Unix())
	endpoint := serviceName + "." + namespace + ".svc.cluster.local:" + port
	// endpoint := "play.min.io"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// log.Printf("%#v\n", minioClient) // minioClient is now setup

	createBucket(minioClient, bucketName, "us-east-1")
	if mode == "liveness" {
		createContent(fileName, 2097152)
		livenessCheck(minioClient, bucketName, fileName)
	} else {
		fileSize := int(capacity * 1024 * 1024 * 1024)
		createContent(fileName, fileSize)
		fmt.Println("Load generation may take some time")
		loadBucket(minioClient, bucketName, fileName)
	}

}

func createBucket(client *minio.Client, bucketName string, region string) {
	count := 0
	var err error
	for count == 0 || (err != nil && count < 5) {
		err = client.MakeBucket(bucketName, region)
		fmt.Println(err)
		count++
		time.Sleep(time.Duration(livenessCheckInterval) * time.Second)
	}
	if err != nil {
		return
	}
	fmt.Println("Successfully Created bucket with name: ", bucketName)

}

func createContent(fileName string, fileSize int) {
	bigBuff := make([]byte, fileSize)
	ioutil.WriteFile(fileName, bigBuff, 0666)
}

func loadBucket(client *minio.Client, bucketName string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return err
	}
	n, err := client.PutObject(bucketName, "test-object", file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/text"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", n)
	return nil
}

func unloadBucket(client *minio.Client, bucketName string) error {
	err := client.RemoveObject(bucketName, "test-object")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func livenessCheck(client *minio.Client, bucketName string, fileName string) {
	count := 0
	for {
		if err := loadBucket(client, bucketName, fileName); err != nil {
			count++
			fmt.Println("Liveness failed, Retrying...")
			if count > livenessRetryCount {
				fmt.Println("Liveness Failed")
				os.Exit(1)
			}
			continue
		}
		unloadBucket(client, bucketName)
		time.Sleep(time.Duration(livenessCheckInterval) * time.Second)
		fmt.Println("Liveness Running")
	}
}

func initVars() error {
	var err error
	mode, err = getEnv("MODE")
	if err != nil {
		return err
	}
	port, err = getEnv("PORT")
	if err != nil {
		return err
	}
	accessKeyID, err = getEnv("ACCESS_KEY")
	if err != nil {
		return err
	}
	secretAccessKey, err = getEnv("SECRET_KEY")
	if err != nil {
		return err
	}
	serviceName, err = getEnv("SERVICE_NAME")
	if err != nil {
		return err
	}
	namespace, err = getEnv("NAMESPACE")
	if err != nil {
		return err
	}
	tmpLivenessCheckInterval, err := getEnv("LIVENESS_CHECK_INTERVAL")
	if err != nil {
		return err
	}
	livenessCheckInterval, err = strconv.Atoi(tmpLivenessCheckInterval)
	if err != nil {
		return err
	}
	tmpLivenessRetryCount, err := getEnv("LIVENESS_RETRY_COUNT")
	if err != nil {
		return err
	}
	livenessRetryCount, err = strconv.Atoi(tmpLivenessRetryCount)
	if err != nil {
		return err
	}
	tmpcapacity, err := getEnv("CAPACITY")
	if err != nil {
		return err
	}
	capacity, err = strconv.Atoi(tmpcapacity)
	if err != nil {
		return err
	}

	return nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("ENV %s not found", key)
	}
	return value, nil
}
