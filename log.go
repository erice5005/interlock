package interlock

import "time"

type ConnectionLog struct {
	History map[time.Time]LogEntry
}

type LogEntry struct {
	Stamp   time.Time
	Content DataFrame
}

func NewLog() ConnectionLog {
	return ConnectionLog{
		History: make(map[time.Time]LogEntry),
	}
}

func (cl ConnectionLog) Add(df DataFrame) {
	cl.History[time.Now()] = LogEntry{
		Stamp:   time.Now(),
		Content: df,
	}
}
