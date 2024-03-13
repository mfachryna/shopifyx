package handler

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/Croazt/shopifyx/utils/response"
	apierror "github.com/Croazt/shopifyx/utils/response/error"
	"github.com/Croazt/shopifyx/utils/validation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/valyala/fasthttp"
)

type ImageHandler struct {
	v *validator.Validate
}

func NewImageHandler(v *validator.Validate) *ImageHandler {
	return &ImageHandler{
		v: v,
	}
}

func (im *ImageHandler) Store(ctx *fasthttp.RequestCtx) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error("Failed to read file from form", http.StatusBadRequest)
		return
	}

	if err := validation.ValidateFile(im.v, fileHeader); err != nil {
		response.Error(ctx, apierror.CustomError(http.StatusBadRequest, err.Error()))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(ctx, apierror.CustomServerError("failed to open file"))
	}
	defer file.Close()

	imageData := bytes.Buffer{}
	_, err = imageData.ReadFrom(file)
	if err != nil {
		response.Error(ctx, apierror.CustomServerError("failed to generate access token"))
		return
	}

	imageUrl, err := UploadImageToS3(fileHeader.Filename, imageData.Bytes())
	if err != nil {
		log.Printf("Failed to upload image to S3: %v", err)
		response.Error(ctx, apierror.CustomServerError("failed to upload image"))
		return
	}

	response.GenerateResponse(ctx, 200, struct {
		ImageUrl string `json:"imageUrl"`
	}{
		ImageUrl: imageUrl,
	})
}

func UploadImageToS3(objectKey string, imageData []byte) (string, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
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
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)

	// Upload image to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(imageData),
	})
	if err != nil {
		return "", err
	}

	// Generate public URL for the uploaded image
	imageURL := "https://" + bucketName + ".s3.amazonaws.com/" + objectKey

	return imageURL, nil
}
