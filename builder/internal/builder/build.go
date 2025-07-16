package builder

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// RunBuildProcess performs the build, logs output, and uploads files to S3
func RunBuildProcess() {
	// Load environment variables from .env file


	bucketName := "goploy-outputs"
	if bucketName == "" {
		log.Fatal("builder - RunBuildProcess : BUCKET_NAME environment variable not set")
	}


	projectID := os.Getenv("PROJECT_ID")



	if projectID == "" {
		log.Fatal("builder - RunBuildProcess : PROJECT_ID environment variable not set")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("builder - RunBuildProcess : Could not get current directory: %v", err)
	}

	outputDir := filepath.Join(cwd, "output")
	distDir := filepath.Join(outputDir, "dist")

	fmt.Println("builder - RunBuildProcess : Starting Build Process...")
	PublishLog("Build Started...")

	err = RunShellCommand(outputDir, "npm install && npm run build")
	if err != nil {
		log.Fatalf("builder - RunBuildProcess : Build command failed: %v", err)
	}

	fmt.Println("builder - RunBuildProcess : Build Complete")
	PublishLog("Build Complete")

	filesToUpload, err := ListFilesRecursively(distDir)
	if err != nil {
		log.Fatalf("builder - RunBuildProcess : Could not list files in dist/: %v", err)
	}

	PublishLog("Starting to upload...")

	s3 := NewS3Client()

	prefix := filepath.Join("__outputs", projectID)

	for _, f := range filesToUpload {
		relPath, _ := filepath.Rel(distDir, f)
		key := filepath.Join(prefix, relPath)

		PublishLog("uploading " + relPath)

		err := UploadFileToS3(s3, bucketName, key, f)
		if err != nil {
			log.Printf("builder - RunBuildProcess : Failed to upload %s: %v", f, err)
			continue
		}

		PublishLog("uploaded " + relPath)
	}

	PublishLog("Done")
	fmt.Println("builder - RunBuildProcess : Upload Complete.")
}