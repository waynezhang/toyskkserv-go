package googleapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func TransliterateRequest(key string) string {
	u := "http://www.google.com/transliterate?langpair=ja-Hira|ja&text=" + url.QueryEscape(key)
	slog.Info("Fallback to Google", "req", u)

	resp, err := http.Get(u)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response", "err", err)
		return ""
	}

	return paraseResponse(body, key)
}

func paraseResponse(resp []byte, key string) string {
	var data [][]interface{}
	if err := json.Unmarshal([]byte(resp), &data); err != nil {
		slog.Error("Failed to parse response", "resp", string(resp))
		return ""
	}

	if len(data) == 0 {
		return ""
	}

	cdds := data[0]
	if len(cdds) < 2 {
		return ""
	}

	// Only return valid key
	if cdds[0].(string) != key {
		return ""
	}

	ret := ""
	for _, c := range cdds[1].([]interface{}) {
		str := c.(string)
		if str == key {
			continue
		}
		ret = ret + "/" + str
	}

	return ret
}
