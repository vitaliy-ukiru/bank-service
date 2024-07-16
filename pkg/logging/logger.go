package logging

import (
	"context"
	"io"
	"log/slog"
)

type Kind int

type Attr struct {
	attr slog.Attr
}

func newAttr(attr slog.Attr) Attr {
	return Attr{attr: attr}
}

func String(key string, value string) Attr {
	return newAttr(slog.String(key, value))
}

func Int64(key string, value int64) Attr {
	return newAttr(slog.Int64(key, value))
}

func AccountId(accountId int64) Attr {
	return newAttr(slog.Int64("account_id", accountId))
}

func Float64(key string, value float64) Attr {
	return newAttr(slog.Float64(key, value))
}

//
//type Logger interface {
//	Info(op string, msg string, args ...Attr)
//	Error(op string, msg string, err error, args ...Attr)
//}

type Logger struct {
	log *slog.Logger
}

func New(w io.Writer) Logger {
	log := slog.New(slog.NewJSONHandler(
		w,
		nil,
	))
	slog.SetDefault(log)
	return Logger{log: log}

}

func (l Logger) Info(op string, msg string, args ...Attr) {
	l.log.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		msg,
		prepareAttrs(args, slog.String("op", op))...,
	)
}

func (l Logger) Error(op string, msg string, err error, args ...Attr) {
	l.log.LogAttrs(
		context.Background(),
		slog.LevelError,
		msg,
		prepareAttrs(
			args,
			slog.String("op", op),
			slog.Any("error", err),
		)...,
	)
}

func (l Logger) With(attrs ...Attr) Logger {
	return Logger{
		slog.New(l.log.Handler().WithAttrs(prepareAttrs(attrs))),
	}
}

func (l Logger) ToStd() *slog.Logger {
	return l.log
}

func prepareAttrs(attrs []Attr, preparedAttrs ...slog.Attr) []slog.Attr {
	// allocate len(preparedAttrs) + len(attrs),
	// but set as available only prepared
	countOfPrepared := len(preparedAttrs)
	res := make([]slog.Attr, countOfPrepared, countOfPrepared+len(attrs))
	if countOfPrepared != 0 {
		copy(res, preparedAttrs)
	}
	for _, attr := range attrs {
		res = append(res, attr.attr)
	}
	return res
}
