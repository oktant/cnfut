package entities

type SourceDestination struct {
	Source          string `json:"source"  validate:"required"`
	Destination     string `json:"destination"  validate:"required"`
	SourceType      string `json:"sourceType"  validate:"required,oneof=azure s3 google local"`
	DestinationType string `json:"destinationType"  validate:"required,oneof=azure s3 google local"`
	Concurrent      bool   `json:"concurrent"  validate:"omitempty"`
	Region          string `json:"region"  validate:"omitempty"`
}

var SupportedSystems = []string{"azure", "local", "s3"}

const DefaultAwsRegion = "eu-north-1"
