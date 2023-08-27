package api

import (
	"bitmoi/backend/utilities"
	"errors"
	"fmt"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/h2non/filetype"

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
	region              = "ap-northeast-2"
	cloudFront          = "https://cdn.bitmoi.co.kr"
	formFileKey         = "image"
	maxProfileImageSize = 5 * 1024 * 1024
	maxAdImageSize      = 10 * 1024 * 1024
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

func (s *Server) uploadProfileImageToS3(f *multipart.FileHeader, userId string) (string, error) {

	openedFile, err := f.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	buf := make([]byte, 261)
	if _, err := openedFile.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read the file")
	}

	kind, _ := filetype.Match(buf)
	ext := strings.ToLower(kind.Extension)

	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("invalid image format. Only JPEG, PNG, and GIF are allowed")
	}

	if f.Size > maxProfileImageSize {
		return "", fmt.Errorf("image size must be less than %dMB", maxProfileImageSize/1024/1024)
	}

	// buf에 261바이트를 기록하기 위해 진전한 read pointer를 다시 초기화 해야함.
	_, err = openedFile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to reset file pointer")
	}

	contentType := mime.TypeByExtension("." + kind.Extension)
	fileUrl := fmt.Sprintf("%s/%s", cloudFront, userId)

	resp, err := s.s3Uploader.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &userId,
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			_, createErr := s.s3Uploader.PutObject(&s3.PutObjectInput{
				Bucket:      &bucketName,
				Key:         &userId,
				Body:        openedFile,
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
		Key:         &userId,
		Body:        openedFile,
		ContentType: &contentType,
	})

	return fileUrl, updateErr
}

func (s *Server) uploadADImageToS3(f *multipart.FileHeader, userId, location string) (string, error) {

	filePath := fmt.Sprintf("bidding/%s/%s", location, userId)

	openedFile, err := f.Open()
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	buf := make([]byte, 261)
	if _, err := openedFile.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read the file")
	}

	kind, _ := filetype.Match(buf)
	ext := strings.ToLower(kind.Extension)

	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("invalid image format. Only JPEG, PNG, and GIF are allowed")
	}

	if f.Size > maxAdImageSize {
		return "", fmt.Errorf("image size must be less than %dMB", maxAdImageSize/1024/1024)
	}

	// buf에 261바이트를 기록하기 위해 진전한 read pointer를 다시 초기화 해야함.
	_, err = openedFile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to reset file pointer")
	}

	contentType := mime.TypeByExtension("." + kind.Extension)
	fileUrl := fmt.Sprintf("%s/%s", cloudFront, filePath)

	resp, err := s.s3Uploader.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &filePath,
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			_, createErr := s.s3Uploader.PutObject(&s3.PutObjectInput{
				Bucket:      &bucketName,
				Key:         &filePath,
				Body:        openedFile,
				ContentType: &contentType,
			})
			if createErr == nil {
				log.Info().Msgf("%s's ad image of %s page stored in s3 bucket successfully", userId, location)
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
		Key:         &filePath,
		Body:        openedFile,
		ContentType: &contentType,
	})

	return fileUrl, updateErr
}

func (s *Server) deleteObject(key string) {
	input := s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}
	_, err := s.s3Uploader.DeleteObject(&input)
	if err != nil {
		log.Err(err).Msgf("cannot delete object. key: %s", key)
	}
}
