package moderation

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/johnrirwin/flyingforge/internal/models"
)

// AWSDetector calls Rekognition through AWS CLI using byte payloads (no S3 dependency).
type AWSDetector struct {
	region  string
	awsPath string
}

// NewAWSDetector creates a detector that uses ambient AWS credentials/profile.
func NewAWSDetector(ctx context.Context, region string) (*AWSDetector, error) {
	_ = ctx

	awsPath, err := exec.LookPath("aws")
	if err != nil {
		return nil, fmt.Errorf("aws cli not found in PATH")
	}

	return &AWSDetector{
		region:  strings.TrimSpace(region),
		awsPath: awsPath,
	}, nil
}

// DetectModerationLabels calls Rekognition DetectModerationLabels with raw image bytes.
func (d *AWSDetector) DetectModerationLabels(ctx context.Context, imageBytes []byte) ([]models.ModerationLabel, error) {
	if len(imageBytes) == 0 {
		return nil, fmt.Errorf("image bytes are required")
	}

	payload := struct {
		Image struct {
			Bytes string `json:"Bytes"`
		} `json:"Image"`
	}{}
	payload.Image.Bytes = base64.StdEncoding.EncodeToString(imageBytes)

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal rekognition payload: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "rekognition-input-*.json")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	if _, err := tmpFile.Write(payloadBytes); err != nil {
		tmpFile.Close()
		return nil, fmt.Errorf("write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("close temp file: %w", err)
	}

	args := []string{
		"rekognition",
		"detect-moderation-labels",
		"--cli-input-json",
		"file://" + tmpPath,
		"--output",
		"json",
	}
	if d.region != "" {
		args = append(args, "--region", d.region)
	}

	cmd := exec.CommandContext(ctx, d.awsPath, args...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("rekognition cli call failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	var parsed struct {
		ModerationLabels []struct {
			Name       string  `json:"Name"`
			ParentName string  `json:"ParentName"`
			Confidence float64 `json:"Confidence"`
		} `json:"ModerationLabels"`
	}
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("parse rekognition response: %w", err)
	}

	labels := make([]models.ModerationLabel, 0, len(parsed.ModerationLabels))
	for _, label := range parsed.ModerationLabels {
		labels = append(labels, models.ModerationLabel{
			Name:       label.Name,
			ParentName: label.ParentName,
			Confidence: label.Confidence,
		})
	}

	return labels, nil
}
