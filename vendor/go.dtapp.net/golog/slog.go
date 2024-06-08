package golog

import (
	"context"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
)

type SLogFun func() *SLog

type sLogConfig struct {
	showLine               bool              // 显示代码行
	setDefault             bool              // 设置为默认的实例
	setDefaultCtx          bool              // 设置默认上下文
	lumberjackConfig       lumberjack.Logger // 配置lumberjack
	lumberjackConfigStatus bool
}

type SLog struct {
	option         sLogConfig
	logger         *slog.Logger
	jsonHandler    *slog.JSONHandler
	jsonCtxHandler *ContextHandler
}

// NewSlog 创建
func NewSlog(opts ...SLogOption) *SLog {
	sl := &SLog{}
	for _, opt := range opts {
		opt(sl)
	}
	sl.start()
	return sl
}

func (sl *SLog) start() {

	opts := slog.HandlerOptions{
		AddSource: sl.option.showLine, // 输出日志语句的位置信息
		Level:     slog.LevelDebug,    // 设置最低日志等级
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey { // 格式化 key 为 "time" 的属性值
				a.Value = slog.StringValue(a.Value.Time().Format(gotime.DateTimeFormat))
				//return slog.Attr{}
			}
			return a
		},
	}

	// json格式输出
	var mw io.Writer
	if sl.option.lumberjackConfigStatus {
		// 同时控制台和文件输出日志
		mw = io.MultiWriter(os.Stdout, &sl.option.lumberjackConfig)
	} else {
		// 只在文件输出日志
		mw = io.MultiWriter(os.Stdout)
	}

	// 控制台输出
	sl.jsonHandler = slog.NewJSONHandler(mw, &opts)

	// 设置默认上下文
	if sl.option.setDefaultCtx {
		sl.jsonCtxHandler = &ContextHandler{sl.jsonHandler}
		sl.logger = slog.New(sl.jsonCtxHandler)
	} else {
		sl.logger = slog.New(sl.jsonHandler)
	}

	// 将这个 slog 对象设置为默认的实例
	if sl.option.setDefault {
		slog.SetDefault(sl.logger)
	}

}

// WithLogger 跟踪编号
func (sl *SLog) WithLogger() (logger *slog.Logger) {
	if sl.option.setDefaultCtx {
		logger = slog.New(sl.jsonCtxHandler)
	} else {
		logger = slog.New(sl.jsonHandler)
	}
	return logger
}

// WithTraceId 跟踪编号
func (sl *SLog) WithTraceId(ctx context.Context) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String(gotrace_id.TraceIdKey, gotrace_id.GetTraceIdContext(ctx)),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceID 跟踪编号
func (sl *SLog) WithTraceID(ctx context.Context) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String(gotrace_id.TraceIdKey, gotrace_id.GetTraceIdContext(ctx)),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceIdStr 跟踪编号
func (sl *SLog) WithTraceIdStr(traceID string) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String(gotrace_id.TraceIdKey, traceID),
	})
	logger := slog.New(jsonHandler)
	return logger
}

// WithTraceIDStr 跟踪编号
func (sl *SLog) WithTraceIDStr(traceID string) *slog.Logger {
	jsonHandler := sl.jsonHandler.WithAttrs([]slog.Attr{
		slog.String(gotrace_id.TraceIdKey, traceID),
	})
	logger := slog.New(jsonHandler)
	return logger
}
