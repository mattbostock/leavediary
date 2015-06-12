package recovery

import (
	"net/http"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/mattbostock/leavediary/handler"
)

type middleware struct {
	logger    *logrus.Logger
	stackAll  bool
	stackSize int
}

// New returns a new instance of middleware
func New(log *logrus.Logger) *middleware {
	return &middleware{
		logger:    log,
		stackAll:  false,
		stackSize: 1024 * 8,
	}
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			handler.ShowError(w, http.StatusInternalServerError, "")
			stack := make([]byte, m.stackSize)
			stack = stack[:runtime.Stack(stack, m.stackAll)]

			f := "PANIC: %s\n%s"
			m.logger.Errorf(f, err, stack)
		}
	}()

	next(w, r)
}
