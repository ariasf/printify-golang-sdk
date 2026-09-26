package domain

import "time"

// Upload is an uploaded image asset.
type Upload struct {
	ID         UploadID
	FileName   string
	Height     int
	Width      int
	Size       int
	MimeType   string
	PreviewURL string
	UploadTime time.Time
}

// UploadImage is the payload to upload an image, by URL or base64 contents.
// Exactly one of URL or Contents must be set.
type UploadImage struct {
	FileName string
	URL      string
	Contents string
}
