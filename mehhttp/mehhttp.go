// Package mehhttp provides functionality for handling and responding errors.
package mehhttp

import (
	"github.com/lefinal/meh"
	"github.com/lefinal/meh/mehlog"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

var (
	// httpStatusCodeMapper is the HTTPStatusCodeMapper to use in
	// LogAndRespondError. It returns http.StatusInternalServerError per default.
	httpStatusCodeMapper HTTPStatusCodeMapper = func(code meh.Code) int {
		return http.StatusInternalServerError
	}
	// httpStatusCodeMapperMutex locks httpStatusCodeMapper.
	httpStatusCodeMapperMutex sync.RWMutex
)

// HTTPStatusCodeMapper maps a meh.Code to HTTP status code. Used for example in
// LogAndRespondError.
type HTTPStatusCodeMapper func(code meh.Code) int

// SetHTTPStatusCodeMapping sets the mapping of meh.Code to HTTP status code
// that is used in LogAndRespondError.
func SetHTTPStatusCodeMapping(mapping HTTPStatusCodeMapper) {
	httpStatusCodeMapperMutex.Lock()
	defer httpStatusCodeMapperMutex.Unlock()
	httpStatusCodeMapper = mapping
}

// HTTPStatusCode retrieves the HTTP status code for the given error.
func HTTPStatusCode(e error) int {
	httpStatusCodeMapperMutex.RLock()
	defer httpStatusCodeMapperMutex.RUnlock()
	return httpStatusCodeMapper(meh.ErrorCode(e))
}

const (
	// ErrCommunication is used for all problems regarding client communication. As
	// communication is unstable by nature, this should not be reported as classic
	// meh.ErrInternal.
	ErrCommunication meh.Code = "mehhttp-communication"
	// ErrServiceNotReachable is used for problems with requesting third-party
	// services.
	ErrServiceNotReachable meh.Code = "mehhttp-service-not-reachable"
)

// LogAndRespondError logs the given meh.Error and responds using the status
// code mapping set via SetHTTPStatusCodeMapping. The responded message will
// always be empty.
func LogAndRespondError(logger *zap.Logger, w http.ResponseWriter, r *http.Request, e error) {
	// Add request details.
	e = meh.ApplyDetails(e, meh.Details{
		"http_req_url":         r.URL.String(),
		"http_req_host":        r.Host,
		"http_req_method":      r.Method,
		"http_req_user_agent":  r.UserAgent(),
		"http_req_remote_addr": r.RemoteAddr,
	})
	mehlog.Log(logger, e)
	httpStatus := HTTPStatusCode(e)
	err := respondHTTP(w, "", httpStatus)
	if err != nil {
		mehlog.Log(logger, meh.Wrap(err, "respond http", meh.Details{
			"status": httpStatus,
		}))
		return
	}
}

// respondHTTP responds the given message with the status to the
// http.ResponseWriter.
func respondHTTP(w http.ResponseWriter, message string, status int) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	_, err := w.Write([]byte(message))
	if err != nil {
		return &meh.Error{
			Code:       ErrCommunication,
			WrappedErr: err,
			Message:    "write",
		}
	}
	return nil
}
