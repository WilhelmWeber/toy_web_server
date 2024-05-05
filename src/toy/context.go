package toy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/WilhelmWeber/toy_web_server/src/final"
)

type ToyContext struct {
	req_header map[string]string
	req_body   string
	//req_file   []byte //multipart-formdataのファイル
	res_header map[string]string
	params     map[string]string
	query      map[string]string
	form       map[string]string
	ctx_params map[string]string
	//TODO: Responseの追加
	response *ToyReponse
}

func NewContext() *ToyContext {
	return &ToyContext{
		req_header: make(map[string]string),
		req_body:   "",
		res_header: make(map[string]string),
		params:     make(map[string]string),
		query:      make(map[string]string),
		form:       make(map[string]string),
		ctx_params: map[string]string{},
		response:   nil,
	}
}

func (t *ToyContext) HasReponse() bool {
	return t.response != nil
}

func (t *ToyContext) Text(statusCode int, text string) {
	t.SetResponseHeader("Content-Type", final.TEXT)
	t.response = &ToyReponse{
		HttpStatus:  final.STATUSES[statusCode],
		ReponseBody: text,
	}
}

func (t *ToyContext) JSON(statusCode int, jsonMap map[string]interface{}) {
	json, err := json.Marshal(jsonMap)
	if err != nil {
		t.response = &ToyReponse{
			HttpStatus:  final.STATUSES[http.StatusInternalServerError],
			ReponseBody: err.Error(),
		}
		return
	}
	t.SetResponseHeader("Content-Type", final.APPLICATION_JSON)
	t.response = &ToyReponse{
		HttpStatus:  final.STATUSES[statusCode],
		ReponseBody: string(json),
	}
}

func (t *ToyContext) File(statusCode int, file_path string) {
	file, err := os.Open(file_path)
	if err != nil {
		t.SetResponseHeader("Content-Type", final.TEXT_HTML_UTF8)
		t.response = &ToyReponse{
			HttpStatus:  final.STATUSES[http.StatusInternalServerError],
			ReponseBody: err.Error(),
		}
		return
	}

	defer file.Close()

	var buf bytes.Buffer
	_, err2 := io.Copy(&buf, file)
	if err2 != nil {
		t.SetResponseHeader("Content-Type", final.TEXT_HTML_UTF8)
		t.response = &ToyReponse{
			HttpStatus:  final.STATUSES[http.StatusInternalServerError],
			ReponseBody: err2.Error(),
		}
		return
	}

	mime := http.DetectContentType(buf.Bytes())
	t.SetResponseHeader("Content-Type", mime)
	t.response = &ToyReponse{
		HttpStatus:  final.STATUSES[statusCode],
		ReponseBody: buf.String(),
	}

}

func (t *ToyContext) GetResponse(http_ver string) (string, string) {
	if t.response != nil {
		res_header := t.res_header
		var formated_header string
		for k, v := range res_header {
			formated_header += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		return t.response.GetResponse(http_ver, formated_header)
	} else {
		status := final.STATUSES[http.StatusInternalServerError]
		return fmt.Sprintf("%s %s\r\n%s\r\n%s", http_ver, status, "Content-Type: text\r\n", "Error:Response Not Found!"), status
	}
}

func (t *ToyContext) Default404(http_ver string) (string, string) {
	status := final.STATUSES[http.StatusNotFound]
	return fmt.Sprintf("%s %s\r\n%s\r\n%s", http_ver, status, "Content-Type: text\r\n", "Error:Resource Not Found!"), status
}

func (t *ToyContext) Default400(http_ver string) (string, string) {
	status := final.STATUSES[http.StatusBadRequest]
	return fmt.Sprintf("%s %s\r\n%s\r\n%s", http_ver, status, "Content-Type: text\r\n", "Error:Method Not Found!"), status
}

func (t *ToyContext) GetReqHeader(key string) string {
	return t.req_header[key]
}

func (t *ToyContext) SetReqHeader(key string, value string) {
	t.req_header[key] = value
}

func (t *ToyContext) SetReqBody(body string) {
	t.req_body = body
}

func (t *ToyContext) GetReqBody() string {
	return t.req_body
}

func (t *ToyContext) GetParams(key string) string {
	return t.ctx_params[key]
}

func (t *ToyContext) SetParams(key string, value string) {
	t.ctx_params[key] = value
}

func (t *ToyContext) SetResponseHeader(key string, value string) {
	t.res_header[key] = value
}

func (t *ToyContext) RemoveResponseHeader(key string) {
	delete(t.res_header, key)
}

func (t *ToyContext) SetUrlParams(key string, value string) {
	t.params[key] = value
}

func (t *ToyContext) GetUrlParams(key string) string {
	return t.params[key]
}

func (t *ToyContext) SetQueryParams(key string, value string) {
	t.query[key] = value
}

func (t *ToyContext) GetQueryParams(key string) string {
	return t.query[key]
}

func (t *ToyContext) SetFormParams(key string, value string) {
	t.form[key] = value
}

func (t *ToyContext) GetFormParams(key string) string {
	return t.form[key]
}
