package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	storage "google.golang.org/api/storage/v1"
)

const (
	// This can be changed to any valid object name.
	objectName = "test-file"
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = storage.DevstorageFullControlScope
)

func fatalf(service *storage.Service, errorMessage string, args ...interface{}) {
	log.Fatalf("Dying with error:\n"+errorMessage, args...)
}

func main() {

	email := os.Getenv("GCS_EMAIL")
	if email == "" {
		log.Fatalf("$GCS_EMAIL is not provided!")
		os.Exit(1)
	}

	privateKey := os.Getenv("GCS_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatalf("$GCS_PRIVATE_KEY is not provided!")
		os.Exit(1)
	}

	privateKey = strings.Replace(privateKey, "\\n", "\n", -1)

	bucketName := os.Getenv("GCS_BUCKET")
	if bucketName == "" {
		log.Fatalf("$GCS_BUCKET is not provided!")
		os.Exit(1)
	}

	projectID := os.Getenv("GCS_PROJECT_ID")
	if projectID == "" {
		log.Fatalf("$GCS_PROJECT_ID is not provided!")
		os.Exit(1)
	}

	targetFolder := os.Getenv("GCS_FOLDER")
	if targetFolder == "" {
		log.Fatalf("$GCS_FOLDER is not provided!")
		os.Exit(1)
	}

	fileName := os.Getenv("BITRISE_IPA_PATH")
	if fileName == "" {
		log.Fatalf("$BITRISE_IPA_PATH is not provided!")
		os.Exit(1)
	}

	branch := os.Getenv("BITRISE_GIT_BRANCH")
	if branch == "" {
		log.Fatalf("$BITRISE_GIT_BRANCH is not provided!")
		os.Exit(1)
	}

	conf := &jwt.Config{
		Email: email,
		PrivateKey: []byte(privateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/devstorage.read_write",
		},
		TokenURL: google.JWTTokenURL,
	}

	client := conf.Client(oauth2.NoContext)

	// if err != nil {
	// 	log.Fatalf("Unable to get default client: %v", err)
	// }
	service, err := storage.New(client)
	if err != nil {
		log.Fatalf("Unable to create storage service: %v", err)
	}

	// If the bucket already exists and the user has access, warn the user, but don't try to create it.
	if _, err := service.Buckets.Get(bucketName).Do(); err == nil {
		fmt.Printf("Bucket %s already exists - skipping buckets.insert call.", bucketName)
	} else {
		// Create a bucket.
		if res, err := service.Buckets.Insert(projectID, &storage.Bucket{Name: bucketName}).Do(); err == nil {
			fmt.Printf("Created bucket %v at location %v\n\n", res.Name, res.SelfLink)
		} else {
			fatalf(service, "Failed creating bucket %s: %v", bucketName, err)
		}
	}

	// Insert an object into a bucket.
	object := &storage.Object{Name: fmt.Sprintf("%v/%v/build.ipa", targetFolder, branch)}
	file, err := os.Open(fileName)
	if err != nil {
		fatalf(service, "Error opening %q: %v", fileName, err)
	}
	if res, err := service.Objects.Insert(bucketName, object).Media(file).Do(); err == nil {
		fmt.Printf("Created object %v at location %v\n\n", res.Name, res.SelfLink)
	} else {
		fatalf(service, "Objects.Insert failed: %v", err)
	}
}