package entities

type SourceDestination struct {
	Source               string `json:"source"  validate:"required"`
	Destination          string `json:"destination"  validate:"required"`
	SourceType           string `json:"sourceType"  validate:"required,oneof=azure s3 google local"`
	DestinationType      string `json:"destinationType"  validate:"required,oneof=azure s3 google local"`
	Concurrent           bool   `json:"concurrent"  validate:"omitempty"`
	Region               string `json:"region"  validate:"omitempty"`
	Bucket               string `json:"bucket"  validate:"omitempty"`
	GoogleCredentialPath string `json:"googleCredentialPath"  validate:"omitempty"`
	S3AccessKeyId        string `json:"s3AccessKeyId'"  validate:"omitempty"`
	S3SecretAccessKey    string `json:"s3SecretAccessKey'"  validate:"omitempty"`
	Endpoint             string `json:"endpoint'"  validate:"omitempty"`
}

const DefaultAwsRegion = "eu-north-1"
