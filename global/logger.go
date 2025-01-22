package global

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"log/slog"
	"runtime"
	"time"
)

func InitLogger() {
	file := fmt.Sprintf("%s/%s", Conf.Log.Path, Conf.Log.Name)

	writer, err := rotatelogs.New(
		file+"_%Y%m%d.log",
		rotatelogs.WithLinkName(file+".log"),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Fatal("Init logger error: ", err)
	}

	// 全局替换
	slog.SetDefault(slog.New(&Handler{
		w:     writer,
		level: getLogLevel(Conf.Log.Level),
	}))
}

func getLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type Handler struct {
	w     *rotatelogs.RotateLogs
	level slog.Level
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return h.level <= level
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	// 构建自定义日志格式
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	logLevel := r.Level.String()
	msg := r.Message

	// 获取调用者信息
	pc, _, line, ok := runtime.Caller(3) // 跳过的调用堆栈深度，可根据需要调整
	var callerInfo string
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		callerInfo = fmt.Sprintf("%s:%d", funcName, line)
	} else {
		callerInfo = "unknown caller"
	}

	// 读取键值对
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf("{%s: %v} ", a.Key, a.Value)
		return true
	})

	// 输出日志
	logLine := fmt.Sprintf("[%s] (%s): 『%s』 %s %s\n", logLevel, timeStr, callerInfo, msg, attrs)
	_, err := h.w.Write([]byte(logLine))
	return err
}

func (h *Handler) WithAttrs(_ []slog.Attr) slog.Handler {
	// 可在此实现上下文附加逻辑（示例简单返回自己）
	return h
}

func (h *Handler) WithGroup(_ string) slog.Handler {
	// 可在此实现分组逻辑（示例简单返回自己）
	return h
}
