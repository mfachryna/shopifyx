package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
)

type ImageHandler struct {
	v *validator.Validate
}

func NewImageHandler(v *validator.Validate) *ImageHandler {
	return &ImageHandler{
		v: v,
	}
}

func (im *ImageHandler) Store(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		response.Error(w, apierror.CustomError(http.StatusBadRequest, err.Error()))
		return
	}
	defer file.Close()

	// Validate file MIME type
	if err := validateImageFileType(fileHeader); err != nil {
		response.Error(w, apierror.CustomError(http.StatusBadRequest, err.Error()))
		return
	}
	// Validate file size
	if fileHeader.Size > (2 * 1024 * 1024) { // 2 MB
		response.Error(w, apierror.CustomError(http.StatusBadRequest, "File size exceeds the limit (2MB)"))
		return
	}

	imageUrl, err := UploadImageToS3(fileHeader.Filename, file)
	if err != nil {
		fmt.Printf("Failed to upload image to S3: %v", err)
		response.Error(w, apierror.CustomServerError("failed to upload image, server error"))
		return
	}

	response.GenerateResponse(w, 200, struct {
		ImageUrl string `json:"imageUrl"`
	}{
		ImageUrl: imageUrl,
	})
}

func UploadImageToS3(fileName string, image multipart.File) (string, error) {

	bucketName := os.Getenv("S3_BASE_URL")
	s3Id := os.Getenv("S3_ID")
	s3SecretKey := os.Getenv("S3_SECRET_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			s3Id,
			s3SecretKey,
			"",
		),
	})

	fileName = generateRandomString(10) + time.Now().Format("20060102150405") + "-" + fileName
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   image,
	})

	if err != nil {
		return "", err
	}

	imageURL := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucketName, "ap-southeast-1", fileName)
	return imageURL, nil
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func validateImageFileType(fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, ".")+1:])
	if !(ext == "jpg" || ext == "jpeg") {
		return fmt.Errorf("File format must be JPG or JPEG")
	}
	return nil
}
