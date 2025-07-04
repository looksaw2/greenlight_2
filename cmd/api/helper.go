package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type envelope map[string]any

func readURLID(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("Parse id is failed,Please input the correct id number")
	}
	if id <= 0 {
		return 0, errors.New("id must be positive")
	}
	return id, nil
}

func (app *Application) writeJSON(w http.ResponseWriter,
	status int,
	data any,
	headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		//输入的json的格式不对
		var syntaxError *json.SyntaxError
		// 转换时间错，和dst出错
		var unMarshalTypeError *json.UnmarshalTypeError
		//dst本身有问题
		var inValidUnmarshalError *json.InvalidUnmarshalError
		//超出最大的Bytes限制
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly-formed")
		case errors.As(err, &unMarshalTypeError):
			if unMarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unMarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON")
		case errors.Is(err, io.EOF):
			return fmt.Errorf("JSON shouldn't be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json : unknown field")
			return fmt.Errorf("body contains unknown field name : %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d", maxBytesError.Limit)
		case errors.As(err, &inValidUnmarshalError):
			panic("dst is wrong")
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body must contain only one json file")
	}
	return nil
}
