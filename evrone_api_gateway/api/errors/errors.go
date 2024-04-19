package errors

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrAuthDataNotFound = ErrAuthData(errors.New("failed to fetch authentication data"))
	ErrNotFound         = &ErrResponse{HTTPStatusCode: 404, ErrorText: "resource not found."}
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	AppCode   int64             `json:"code,omitempty"`  // application-specific error code
	ErrorText string            `json:"error,omitempty"` // application-level error message, for debugging
	Errors    map[string]string `json:"errors,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequestRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		ErrorText:      err.Error(),
	}
}

func ErrInvalidArgumentRender(err error, errors map[string]string) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		ErrorText:      err.Error(),
		Errors:         errors,
	}
}

func IsNotFound(err error) bool {
	st, ok := status.FromError(err)
	if !ok {
		return false
	}

	return st.Code() == codes.NotFound
}

func ErrAuthData(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 401,
		ErrorText:      err.Error(),
	}
}

func Error(err error) render.Renderer {
	st, ok := status.FromError(err)
	if !ok {
		return ErrInvalidRequestRender(err)
	}

	switch st.Code() {
	// error bad request
	case codes.OK:
		return nil
	// error unavailable
	case codes.AlreadyExists:
		return ErrInvalidRequestRender(errors.New(st.Message()))
	// error unavailable
	case codes.Unavailable:
		return ErrInvalidRequestRender(errors.New("Service is unavailable. Try again soon."))
	// error not found
	case codes.NotFound:
		return ErrNotFound
	// error validation errors
	case codes.InvalidArgument:
		details := ErrorDetails(st)
		return ErrInvalidArgumentRender(ErrInvalidArgument, details)
	default:
		var errStr string
		for _, detail := range st.Details() {
			if errorInfo, ok := detail.(*epb.ErrorInfo); ok {
				errStr += " " + errorInfo.Reason
			}
		}
		if len(errStr) != 0 {
			return ErrInvalidRequestRender(errors.New(errStr))
		}
		return ErrInvalidRequestRender(err)
	}
}

func ErrorDetails(st *status.Status) map[string]string {
	var data = make(map[string]string)
	for _, detail := range st.Details() {
		if badRequest, ok := detail.(*epb.BadRequest); ok {
			for _, violation := range badRequest.GetFieldViolations() {
				data[violation.GetField()] = violation.GetDescription()
			}
		}
	}
	return data
}
