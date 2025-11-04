package fixture

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Beigelman/nossas-despesas/internal/pkg/db"
)

func ExecuteSQLFiles(db *db.Client, files []string) error {
	ctx := context.Background()

	if db == nil {
		return errors.New("client is nil")
	}

	var errs error
	for _, file := range files {
		filePath, err := filepath.Abs(file)
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("getting absolute file (%s) path: %w", file, err))
			continue
		}

		content, err := os.ReadFile(filepath.Clean(filePath))
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("getting absolute file (%s) path: %w", file, err))
			continue
		}

		_, err = db.Conn().ExecContext(ctx, string(content))
		if err != nil {
			errs = errors.Join(errs, fmt.Errorf("getting absolute file (%s) path: %w", file, err))
			continue
		}
	}

	return errs
}
