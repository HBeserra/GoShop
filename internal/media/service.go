package media

import (
	"context"
	"github.com/HBeserra/GoShop/domain"
	"github.com/google/uuid"
	"io"
)

type ContentType = string

type Config struct {
	MediaTypes    []ContentType                   `env:"MEDIA_TYPES" envDefault:"image/png,image/jpeg"`
	MaxSize       domain.ByteSize                 `env:"MAX_SIZE" envDefault:"2mb"`
	Compress      bool                            `env:"COMPRESS" envDefault:"true"`
	CompressLevel int                             `env:"COMPRESS_LEVEL" envDefault:"9"`
	MaxUploadSize domain.ByteSize                 `env:"MAX_UPLOAD_SIZE" envDefault:"5mb"`
	MaxSizeByType map[ContentType]domain.ByteSize `env:"MAX_SIZE_BY_TYPE" envDefault:"{}"`
}

type Service struct {
	config     *Config
	auth       AuthService
	storage    FileStorage
	repo       MediaRepository
	compressor map[ContentType]FileCompressor
}

//go:generate mockgen -source=service.go -destination mock_test.go --package  media_test
type FileStorage interface {
	Save(file io.Reader, filename string) (string, error)
	Get(filename string) ([]byte, error)
	GetURL(filename string) (string, error)
	Delete(filename string) error
}

type MediaRepository interface {
	Save(ctx context.Context, namespace string, media *domain.Media) (uuid.UUID, error)
	GetByID(ctx context.Context, namespace string, id uuid.UUID) (*domain.Media, error)
	Delete(ctx context.Context, namespace string, id uuid.UUID) error
}

type FileCompressor interface {
	Compress(file io.Reader, level int) ([]byte, error)
}

type AuthService interface {
	// GetUserID retrieves the unique identifier (UUID) of the user from the provided context. Returns an error if retrieval fails.
	GetUserID(ctx context.Context) (uuid.UUID, error)
	// CheckPermissions verifies if a user has the required permissions for a specified namespace and action(s). It returns a boolean indicating access and an error if the operation fails.
	CheckPermissions(ctx context.Context, userID uuid.UUID, namespace string, permission ...string) (bool, error)
}

func NewService(config *Config, storage FileStorage, repo MediaRepository, auth AuthService, compressors map[string]FileCompressor) *Service {
	return &Service{
		config:     config,
		storage:    storage,
		repo:       repo,
		auth:       auth,
		compressor: compressors,
	}
}
