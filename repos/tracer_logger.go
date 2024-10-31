package repos

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func prettyPrintSQL(sql string) []string {
	return strings.Split(sql, "\n")
}

type LoggingQueryTracer struct {
	logger *slog.Logger
	start  int64
	end    int64
}

func NewLoggingQueryTracer(logger *slog.Logger) *LoggingQueryTracer {
	return &LoggingQueryTracer{logger: logger, start: time.Now().UnixMilli(), end: 0}
}

func (l *LoggingQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	l.logger.Info("[query start]", slog.Any("args", data.Args))
	for _, line := range prettyPrintSQL(data.SQL) {
		l.logger.Info(line)
	}
	return ctx
}

func (l *LoggingQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	l.end = time.Now().UnixMilli()

	if data.Err != nil {
		l.logger.
			Error("[query end]",
				slog.String("error", data.Err.Error()),
				slog.String("command_tag", data.CommandTag.String()),
			)
		return
	}

	l.logger.
		Info("[query end]",
			slog.String("time", strconv.FormatInt(l.end-l.start, 10)+"ms"),
			slog.String("command_tag", data.CommandTag.String()),
		)
}

type MultiQueryTracer struct {
	Tracers []pgx.QueryTracer
}

func NewMultiQueryTracer(tracers ...pgx.QueryTracer) *MultiQueryTracer {
	return &MultiQueryTracer{Tracers: tracers}
}

func (m *MultiQueryTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, t := range m.Tracers {
		ctx = t.TraceQueryStart(ctx, conn, data)
	}
	return ctx
}

func (m *MultiQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, t := range m.Tracers {
		t.TraceQueryEnd(ctx, conn, data)
	}
}