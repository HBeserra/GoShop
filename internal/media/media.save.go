package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/pkg/observability"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"slices"
	"time"
)

func (s *Service) Save(
	ctx context.Context,
	namespace string,
	file io.ReadCloser,
	filename,
	contentType string,
) (uuid.UUID, error) {

	ctx, span := observability.StartSpan(ctx, "media.Save")
	defer span.End()

	defer file.Close()

	userID, err := s.auth.GetUserID(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", domain.ErrUnauthorized, err)
	}

	perm, err := s.auth.CheckPermissions(ctx, userID, namespace, "media:create")
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", domain.ErrUnauthorized, err)
	}

	if !perm {
		return uuid.Nil, domain.ErrUnauthorized
	}

	if !slices.Contains(s.config.MediaTypes, contentType) {
		return uuid.Nil, domain.ErrInvalidMediaType
	}

	var fileData []byte
	if compressor, ok := s.compressor[contentType]; ok {
		fileData, err = compressor.Compress(file, s.config.CompressLevel)
		if err != nil {
			return uuid.Nil, err
		}
	}

	if len(fileData) > int(s.config.MaxSize) {
		return uuid.Nil, domain.ErrFileTooLarge
	}

	reader := bytes.NewReader(fileData)

	url, err := s.storage.Save(reader, filename)
	if err != nil {
		slog.ErrorContext(ctx, "error saving file",
			"error", err,
			"filename", filename,
			"contentType", contentType,
			"size", len(fileData),
			"namespace", namespace,
		)
		return uuid.Nil, err
	}

	s.CreateThumb(ctx, fileData, filename, contentType)

	return s.repo.Save(ctx, namespace, &domain.Media{
		ID:        uuid.New(),
		Namespace: namespace,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Filename:  filename,
		Url:       url,
		Type:      contentType,
		Size:      len(fileData),
	})

}
