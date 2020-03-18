package s3

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func TestUploadMultipart(t *testing.T) {
	S3_REGION := "cn-north-1"
	S3_BUCKET := "agora-finace-backup-test"

	// Upload Files
	// fileName := "../compress/file1.txt"
	// fileName := "/Users/slq/Downloads/Microsoft_Remote_Desktop_Beta.app.zip"
	fileName := "/Users/slq/Downloads/grand-masion-tokyo/Grand-Maison-Tokyo.E11.2019.WEB-DL.1080p.H264.AAC-PTHome.mp4"
	s3Key := "uploadMultipart.txt"
	err := UploadMultipart(S3_REGION, S3_BUCKET, fileName, s3Key)
	if err != nil {
		log.Fatal(err)
	}
}
func TestUpload(t *testing.T) {

	S3_REGION := "cn-north-1"
	S3_BUCKET := "agora-finace-backup-test"

	credentialsChainVerboseErrors := true
	sess, err := session.NewSessionWithOptions(session.Options{
		// Specify profile to load for the session's config
		// Profile: "/Users/slq/.aws/config",

		// Provide SDK Config options, such as Region.
		Config: aws.Config{
			Region:                        aws.String(S3_REGION),
			CredentialsChainVerboseErrors: &credentialsChainVerboseErrors,
			// Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		},

		// Force enable Shared Config support
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Upload Files
	fileName := "../compress/file1.txt"
	err = uploadFile(sess, S3_BUCKET, fileName)
	if err != nil {
		log.Fatal(err)
	}
}
