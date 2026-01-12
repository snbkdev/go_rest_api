package helpers

import (
	"context"
	log "github.com/go-kit/kit/log"
	"time"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next EncryptService
}

func (mw LoggingMiddleware) Encrypt(ctx context.Context, key string, text string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "encrypt", "key", key, "text", text, "output", output, "err", err, "took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.Next.Encrypt(ctx, key, text)
	return
}

func (mw LoggingMiddleware) Decrypt(ctx context.Context, key string, text string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "decrypt", "key", key, "message", text, "output", output, "err", err, "took", time.Since(begin),
		)
	}(time.Now())
	output, err = mw.Next.Decrypt(ctx, key, text)
	return
}