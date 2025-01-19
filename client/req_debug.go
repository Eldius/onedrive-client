package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"reflect"
	"slices"
	"strings"
)

func debugResponse(ctx context.Context, res *http.Response, reqBody []byte) {
	resBody, _ := io.ReadAll(res.Body)
	slog.With("request", map[string]any{
		"status_code": res.StatusCode,
		"body":        string(parseBody(reqBody)),
		"headers":     headerToMap(res.Request.Header),
		"method":      res.Request.Method,
		"url":         res.Request.URL.String(),
		"response": map[string]any{
			"body":    string(parseBody(resBody)),
			"headers": headerToMap(res.Header),
		},
	}).DebugContext(ctx, "externalRequest")
	res.Body = io.NopCloser(bytes.NewReader(resBody))
}

func parseBody(b []byte) []byte {
	var bodyMap map[string]any
	if err := json.Unmarshal(b, &bodyMap); err != nil {
		return b
	}
	bodyMap = parseMap(bodyMap)
	body, err := json.Marshal(bodyMap)
	if err != nil {
		return b
	}
	return body
}

var (
	redactedKeyList = []string{
		"access_token",
		"accesstoken",
		"refresh_token",
		"refreshtoken",
		"token_type",
		"idtoken",
		"authorization",
		"athentication",
	}
)

func parseMap(v map[string]any) map[string]any {
	v = maps.Clone(v)
	for k := range maps.Keys(v) {
		if slices.Contains(redactedKeyList, strings.ToLower(k)) {
			if reflect.TypeOf(v[k]).Kind() == reflect.Map {
				v[k] = parseMap(v[k].(map[string]any))
			}
			v[k] = "***"
		}
	}
	return v
}

func headerToMap(h http.Header) map[string]any {
	m := make(map[string]any)
	for k, v := range h {
		//if strings.Contains(strings.ToLower(k), "authorization") ||
		//	strings.Contains(strings.ToLower(k), "authentication") {
		//	continue
		//}
		m[k] = v
	}
	return parseMap(m)
}
