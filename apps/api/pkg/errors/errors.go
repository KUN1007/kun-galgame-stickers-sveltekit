package errors

import "fmt"

type AppError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func New(code int, message string, statusCode int) *AppError {
	return &AppError{Code: code, Message: message, StatusCode: statusCode}
}

const (
	CodeOK   = 0
	CodeAuth = 205
	CodeBiz  = 233
)

// Domain codes. The three generic codes above cannot tell a caller why a write
// failed -- 233 alone stood for 400, 404, 500 and 503 -- so every outcome a
// client is expected to act on gets its own number. The image service does the
// same thing in the 80001-80015 range; these are the sticker equivalents.
const (
	CodePackNotFound     = 90001
	CodeStickerNotFound  = 90002
	CodeNotOwner         = 90003
	CodePackLimit        = 90004
	CodeStickerLimit     = 90005
	CodeUploadDailyLimit = 90006
	CodeFileTooLarge     = 90007
	CodePackNeedsSticker = 90008
	CodeImageRejected    = 90009
	CodeModeration       = 90010
	CodeImageUnavailable = 90011
	CodeInvalidParams    = 90012
	CodeTagLimit         = 90013
)

func ErrUnauthorized(msg string) *AppError { return New(CodeAuth, msg, 401) }
func ErrForbidden(msg string) *AppError    { return New(CodeAuth, msg, 403) }
func ErrBadRequest(msg string) *AppError   { return New(CodeBiz, msg, 400) }
func ErrNotFound(msg string) *AppError     { return New(CodeBiz, msg, 404) }
func ErrUnavailable(msg string) *AppError  { return New(CodeBiz, msg, 503) }
func ErrInternal(msg string) *AppError     { return New(CodeBiz, msg, 500) }

func ErrPackNotFound() *AppError    { return New(CodePackNotFound, "pack not found", 404) }
func ErrStickerNotFound() *AppError { return New(CodeStickerNotFound, "sticker not found", 404) }
func ErrNotOwner() *AppError        { return New(CodeNotOwner, "not the owner of this pack", 403) }

func ErrInvalidParams(msg string) *AppError { return New(CodeInvalidParams, msg, 400) }

func ErrPackLimit(max int) *AppError {
	return New(CodePackLimit, fmt.Sprintf("at most %d packs per user", max), 400)
}

func ErrStickerLimit(max int) *AppError {
	return New(CodeStickerLimit, fmt.Sprintf("at most %d stickers per pack", max), 400)
}

func ErrUploadDailyLimit(max int) *AppError {
	return New(CodeUploadDailyLimit, fmt.Sprintf("at most %d uploads per day", max), 429)
}

func ErrFileTooLarge(maxBytes int64) *AppError {
	return New(CodeFileTooLarge, fmt.Sprintf("file exceeds %d bytes", maxBytes), 413)
}

func ErrPackNeedsSticker() *AppError {
	return New(CodePackNeedsSticker, "a pack needs at least one sticker to publish", 400)
}

func ErrImageRejected(msg string) *AppError { return New(CodeImageRejected, msg, 400) }

func ErrModeration() *AppError { return New(CodeModeration, "image rejected by moderation", 422) }

func ErrImageUnavailable() *AppError {
	return New(CodeImageUnavailable, "image service unavailable", 503)
}

func ErrTagLimit(max int) *AppError {
	return New(CodeTagLimit, fmt.Sprintf("at most %d tags per pack", max), 400)
}
