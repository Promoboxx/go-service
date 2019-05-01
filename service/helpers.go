package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/promoboxx/go-glitch/glitch"
	"github.com/promoboxx/go-service/alice/middleware/lrw"
)

// ReturnProblem will return a json http problem response
func ReturnProblem(w http.ResponseWriter, detail, code string, status int, innerErr error) (int, []byte) {
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}

	if loggingResponseWriter, ok := w.(*lrw.LoggingResponseWriter); ok {
		loggingResponseWriter.InnerError = innerErr
	}

	by, _ := json.Marshal(prob)
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	return status, by
}

// WriteProblem will write a json http problem response
func WriteProblem(w http.ResponseWriter, detail, code string, status int, innerErr error) error {
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}
	by, err := json.Marshal(prob)
	if err != nil {
		return err
	}

	if lrw, ok := w.(*lrw.LoggingResponseWriter); ok {
		lrw.InnerError = innerErr
	}

	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}

// WriteJSONResponse will write a json response to the htt.ResponseWriter
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	var by []byte
	var err error
	if data != nil {
		by, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}

// Int32PointerFromQueryParam returns a nullable int32 from a query param key
func Int32PointerFromQueryParam(r *http.Request, paramName string) (*int32, error) {
	strValue := r.URL.Query().Get(paramName)
	var intPointer *int32
	if len(strValue) > 0 {
		i, err := strconv.ParseInt(strValue, 10, 32)
		if err != nil {
			return intPointer, err
		}
		i32 := int32(i)
		intPointer = &i32
	}
	return intPointer, nil
}
