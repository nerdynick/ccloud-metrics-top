package widgets

import (
	"fmt"
	"strings"

	"github.com/gizak/termui/v3/widgets"
	log "github.com/sirupsen/logrus"
)

const (
	maxLogRows = 10
)

type LogrusList struct {
	widgets.List
	logLevels []log.Level
}

func (lr *LogrusList) Levels() []log.Level {
	return lr.logLevels
}

func (lr *LogrusList) Fire(entry *log.Entry) error {
	// time := entry.Time
	msg := entry.Message
	lvl := entry.Level
	fields := entry.Data

	strFields := []string{}
	for k, v := range fields {
		strFields = append(strFields, fmt.Sprintf("%s=%v", k, v))
	}

	lr.Rows = append(lr.Rows, fmt.Sprintf("[%s] %s  Fields|%s|", lvl, msg, strings.Join(strFields, ", ")))
	l := len(lr.Rows)
	if l > maxLogRows {
		lr.Rows = lr.Rows[l-maxLogRows : l]
	}
	lr.ScrollBottom()
	return nil
}

func NewLogrusList(lvl ...log.Level) *LogrusList {
	logWidget := LogrusList{
		List:      *widgets.NewList(),
		logLevels: lvl,
	}
	logWidget.Rows = make([]string, 0)

	return &logWidget
}
