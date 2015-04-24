package mockhook

import "github.com/Sirupsen/logrus"

type Mockhook struct {
	Entries []*logrus.Entry
}

func (h *Mockhook) Fire(entry *logrus.Entry) error {
	h.Entries = append(h.Entries, entry)
	return nil
}

func (h *Mockhook) LastEntry() *logrus.Entry {
	return h.Entries[len(h.Entries)-1]
}

func (h *Mockhook) Levels() []logrus.Level {
	// Capture all logrus levels
	return []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
