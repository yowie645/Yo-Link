package slogdiscard

import (
	"context"
	"log/slog"
)

// NewDiscardLogger возвращает новый slog.Logger, который игнорирует все сообщения.
func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

// DiscardHandler реализует slog.Handler, который игнорирует все сообщения.
type DiscardHandler struct{}

// NewDiscardHandler создает новый DiscardHandler.
func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

// Handle игнорирует переданную запись журнала.
func (h *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs возвращает тот же обработчик, так как атрибуты не сохраняются.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup возвращает тот же обработчик, так как группа не сохраняется.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Enabled всегда возвращает false, так как все сообщения игнорируются.
func (h *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Проверка, что DiscardHandler реализует slog.Handler на этапе компиляции
var _ slog.Handler = (*DiscardHandler)(nil)
