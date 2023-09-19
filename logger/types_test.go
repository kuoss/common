package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	testCases := []struct {
		logLevel Level
		want     string
	}{
		{PanicLevel, "panic"},
		{FatalLevel, "fatal"},
		{ErrorLevel, "error"},
		{WarnLevel, "warning"},
		{InfoLevel, "info"},
		{DebugLevel, "debug"},
		{TraceLevel, "trace"},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got := tc.logLevel.String()
			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseLevel(t *testing.T) {
	testCases := []struct {
		levelString string
		want        Level
		wantError   string
	}{
		{"panic", PanicLevel, ``},
		{"PANIC", PanicLevel, ``},
		{"fatal", FatalLevel, ``},
		{"FATAL", FatalLevel, ``},
		{"error", ErrorLevel, ``},
		{"ERROR", ErrorLevel, ``},
		{"warn", WarnLevel, ``},
		{"WARN", WarnLevel, ``},
		{"warning", WarnLevel, ``},
		{"WARNING", WarnLevel, ``},
		{"info", InfoLevel, ``},
		{"INFO", InfoLevel, ``},
		{"debug", DebugLevel, ``},
		{"DEBUG", DebugLevel, ``},
		{"trace", TraceLevel, ``},
		{"TRACE", TraceLevel, ``},
		{"invalid", PanicLevel, `not a valid logrus Level: "invalid"`},
		{"foo", PanicLevel, `not a valid logrus Level: "foo"`},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got, err := ParseLevel(tc.levelString)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, got)
		})
	}
}
