package httpx

import (
	"net/http"
	"net/url"
	"strconv"
)

func ReadIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, err
	}

	return id, nil
}

func ReadInt(qp url.Values, param string, defaultValue int) int {
	s := qp.Get(param)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}
