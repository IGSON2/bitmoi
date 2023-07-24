package api

import (
	"bitmoi/backend/utilities"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/h2non/filetype"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	errCannotUpdateUntil24H = errors.New("profile picture can updated after 24 hours until last modified")
	allowedImageExtensions  = map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"gif":  true,
	}
)

const (
	region       = "ap-northeast-2"
	cloudFront   = "https://d3qqltqqisxphn.cloudfront.net"
	fileKey      = "image"
	userIdKey    = "user_id"
	maxImageSize = 10 * 1024 * 1024
)

var bucketName = "bitmoi-photo"

func NewS3Uploader(c *utilities.Config) (*s3.S3, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(c.S3AccessKey, c.S3SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)

	return svc, nil
}

func (s *Server) uploadImageToS3(c *fiber.Ctx) (string, error) {
	f, err := c.FormFile(fileKey)
	if err != nil {
		return "", fmt.Errorf("cannot get photo image file from context. err: %w", err)
	}
	userId := c.FormValue(userIdKey)
	if userId == "" {
		return "", fmt.Errorf("cannot get user id from context")
	}

	openedF, err := f.Open()
	if err != nil {
		return "", err
	}
	defer openedF.Close()

	buf := make([]byte, 261)
	if _, err := openedF.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read the file")
	}

	kind, _ := filetype.Match(buf)
	ext := strings.ToLower(kind.Extension)
	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("invalid image format. Only JPEG, PNG, and GIF are allowed")
	}

	if f.Size > maxImageSize {
		return "", fmt.Errorf("image size must be less than %dMB", maxImageSize/1024/1024)
	}

	contentType := fmt.Sprintf("image/%s", ext)
	fileKey := fmt.Sprintf("%s.%s", userId, ext)
	fileUrl := fmt.Sprintf("%s/%s", cloudFront, fileKey)

	resp, err := s.s3Uploader.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &fileKey,
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			_, createErr := s.s3Uploader.PutObject(&s3.PutObjectInput{
				Bucket:      &bucketName,
				Key:         &fileKey,
				Body:        openedF,
				ContentType: &contentType,
			})
			if createErr == nil {
				log.Info().Msgf("%s's new photo stored in s3 bucket successfully", userId)
				return fileUrl, createErr
			} else {
				return "", createErr
			}
		} else {
			return "", err
		}
	}

	if time.Since(*resp.LastModified) < 24*time.Hour {
		return "", errCannotUpdateUntil24H
	}

	_, updateErr := s.s3Uploader.PutObject(&s3.PutObjectInput{
		Bucket:      &bucketName,
		Key:         &fileKey,
		Body:        openedF,
		ContentType: &contentType,
	})

	return fileUrl, updateErr
}

func (s *Server) deleteImage(key string) {
	input := s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}
	_, err := s.s3Uploader.DeleteObject(&input)
	if err != nil {
		log.Err(err).Msgf("cannot delete object. key: %s", key)
	}
}

func (s *Server) testUpload(c *fiber.Ctx) error {
	url, err := s.uploadImageToS3(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(url)
}
