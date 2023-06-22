package pkg

import (
	"context"
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

//go:generate easyjson  -disallow_unknown_fields -omit_empty wrapper.go

//easyjson:json
type ErrResponse struct {
	ErrMassage string `json:"message,omitempty"`
}

func DefaultHandlerHTTPError(ctx context.Context, w http.ResponseWriter, err error) {
	errCause := errors.Cause(err)

	code, exist := GetErrorCodeHTTP(errCause)
	if !exist {
		errCause = errors.Wrap(errCause, "Undefined error")
	}

	errResp := ErrResponse{
		ErrMassage: errCause.Error(),
	}

	Response(ctx, w, code, errResp)
}

func NoBody(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func getEasyJSON(someStruct interface{}) ([]byte, error) {
	someStructUpdate, ok := someStruct.(easyjson.Marshaler)
	if !ok {
		return []byte{}, ErrGetEasyJSON
	}

	out, err := easyjson.Marshal(someStructUpdate)
	if !ok {
		return []byte{}, ErrJSONUnexpectedEnd
	}

	return out, err
}

func Response(ctx context.Context, w http.ResponseWriter, statusCode int, someStruct interface{}) {
	out, err := getEasyJSON(someStruct)
	if err != nil {
		DefaultHandlerHTTPError(ctx, w, err)
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)

	w.WriteHeader(statusCode)

	_, err = w.Write(out)
	if err != nil {
		DefaultHandlerHTTPError(ctx, w, err)
		return
	}
}
