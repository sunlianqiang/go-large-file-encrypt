package s3

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/astaxie/beego/logs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadMultipart(s3Region, s3Bucket, uploadFile string) (err error) {
	logs.Info("---------- UploadMultipart START ----------")
	defer logs.Info("---------- UploadMultipart END ----------")

	// The session the S3 Uploader will use
	credentialsChainVerboseErrors := true
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region:                        aws.String(s3Region),
			CredentialsChainVerboseErrors: &credentialsChainVerboseErrors,
			// Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(uploadFile)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", uploadFile, err)
	}
	defer f.Close()

	// Upload the file to S3.
	s3Key := filepath.Base(uploadFile)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(s3Key),
		Body:   f,
	})
	if err != nil {
		err = fmt.Errorf("failed to upload file, %v", err)
		logs.Error(err.Error())
		return
	}
	logs.Debug("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return
}

func UploadFile(s3Region, s3Bucket, uploadFileDir string) {

	credentialsChainVerboseErrors := true
	sess, err := session.NewSessionWithOptions(session.Options{
		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region:                        aws.String(s3Region),
			CredentialsChainVerboseErrors: &credentialsChainVerboseErrors,
			// Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Upload Files
	err = uploadFile(sess, s3Bucket, uploadFileDir)
	if err != nil {
		log.Fatal(err)
		return
	}

}

func uploadFile(session *session.Session, s3Bucket, uploadFileDir string) error {

	upFile, err := os.Open(uploadFileDir)
	if err != nil {
		return err
	}
	defer upFile.Close()

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)

	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s3Bucket),
		Key:                  aws.String(uploadFileDir),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(fileBuffer),
		ContentLength:        aws.Int64(fileSize),
		ContentType:          aws.String(http.DetectContentType(fileBuffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}
