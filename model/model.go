package model

// Metadata object to be stored in S3 buckets
type Metadata struct {
	// Description is the description of the Meme
	Description string `json:"description"`
}
