package media

import (
	"bytes"
	"context"
	"fmt"
	"github.com/HBeserra/GoShop/domain"
	"github.com/HBeserra/GoShop/pkg/observability"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"math"
	"os"
)

// CreateThumb generates a thumbnail for images and videos.
func (s *Service) CreateThumb(
	ctx context.Context,
	file []byte,
	filename,
	contentType string,
) ([]byte, error) {
	if contentType == "image/jpeg" || contentType == "image/png" {
		// Handle image thumbnail generation
		img, _, err := image.Decode(bytes.NewReader(file))
		if err != nil {
			return nil, err
		}

		// Resize the image (simplified logic)
		thumbnail := image.NewRGBA(image.Rect(0, 0, 100, 100)) // Thumbnail size: 100x100
		for y := 0; y < 100; y++ {
			for x := 0; x < 100; x++ {
				thumbnail.Set(x, y, img.At(x*img.Bounds().Dx()/100, y*img.Bounds().Dy()/100))
			}
		}

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, thumbnail, nil)
		if err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	} else if contentType == "video/mp4" {

	}

	return nil, os.ErrInvalid
}

// thumbnail creates a resized image from the reader and writes it to
// the writer. The mimetype determines how the image will be decoded
// and must be either "image/jpeg" or "image/png". The desired width
// of the thumbnail is specified in pixels, and the resulting height
// will be calculated to preserve the aspect ratio.
func thumbnail(r io.Reader, w io.Writer, mimetype string, width int) error {
	var src image.Image
	var err error

	switch mimetype {
	case "image/jpeg":
		src, err = jpeg.Decode(r)
	case "image/png":
		src, err = png.Decode(r)
	}

	if err != nil {
		return err
	}

	ratio := (float64)(src.Bounds().Max.Y) / (float64)(src.Bounds().Max.X)
	height := int(math.Round(float64(width) * ratio))

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	err = jpeg.Encode(w, dst, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createVideoThumb(
	ctx context.Context,
	file []byte,
	filename string,
	contentType string,
	width, height int,
) ([]byte, error) {
	ctx, span := observability.StartSpan(ctx, "media.createVideoThumb")
	defer span.End()

	if contentType != "video/mp4" {
		return nil, domain.ErrInvalidMediaType
	}

	tmpFile, err := os.CreateTemp("", "thumb-*.jpg")
	if err != nil {
		slog.ErrorContext(ctx, "error creating temporary file", "error", err, "filename", tmpFile.Name())
		return nil, err
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			slog.ErrorContext(ctx, "error removing temporary file", "error", err, "filename", tmpFile.Name())
		}
	}(tmpFile.Name())

	_, err = tmpFile.Write(file)
	if err != nil {
		slog.ErrorContext(ctx, "error writing file", "error", err, "filename", tmpFile.Name())
		return nil, err
	}

	err = tmpFile.Close()
	if err != nil {
		slog.ErrorContext(ctx, "error closing file", "error", err, "filename", tmpFile.Name())
		return nil, err
	}

	// Initialize FFmpeg
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(tmpFile.Name()).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 0)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}

	return buf.Bytes(), nil
}
