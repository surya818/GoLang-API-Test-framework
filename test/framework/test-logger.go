package framework

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() error {
	// Defining a custom encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "output", // Only include the "output" field
		LevelKey:      "",       // Exclude level
		TimeKey:       "",       // Exclude timestamp
		NameKey:       "",       // Exclude logger name
		CallerKey:     "",       // Exclude caller information
		StacktraceKey: "",       // Exclude stacktrace
		EncodeTime:    nil,      // No time encoding
	}

	// Create a custom core with the encoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),    // Custom JSON encoder
		zapcore.AddSync(zapcore.Lock(os.Stdout)), // Output to stdout
		zapcore.InfoLevel,                        // Log level
	)
	Logger = zap.New(core)
	defer Logger.Sync()
	return nil
}
