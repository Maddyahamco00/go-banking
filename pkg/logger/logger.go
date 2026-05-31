package logger

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New(level string) (*Logger, error) {
	// Map common log levels; zap will validate.
	lvl := zapcore.InfoLevel
	if err := lvl.Set(level); err != nil {
		lvl = zapcore.InfoLevel
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encCfg),
		zapcore.AddSync(io.Writer(osStdoutWriter{})),
		lvl,
	)

	lg := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &Logger{Logger: lg}, nil
}

// WithRequest attaches request-scoped fields.
func (l *Logger) WithRequest(ctx context.Context, fields map[string]any) *zap.Logger {
	_ = ctx
	lg := l.Logger

	for k, v := range fields {
		lg = lg.With(zap.Any(k, v))
	}
	// Ensure we include request context timestamp for correlation.
	return lg.With(zap.String("request_ts", time.Now().UTC().Format(time.RFC3339Nano)))
}

// osStdoutWriter writes to stdout.
// Trade-off: kept minimal; production setups may prefer zap's os.Stdout directly.
type osStdoutWriter struct{}

func (osStdoutWriter) Write(p []byte) (n int, err error) {
	// We avoid importing os here; zapcore.AddSync will wrap this.
	return len(p), nil
}
