package upload

import (
	"context"
	"os"
)

func (r *Repository) UploadFile(ctx context.Context, file []byte, path string) (string, error) {
	err := os.WriteFile(path, file, 0644)
	if err != nil {
		return "", err
	}
	return "http://dummy.com/" + path, nil
}
