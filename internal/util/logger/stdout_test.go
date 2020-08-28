package logger

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestClientValuesAreDefinedCorrectly(t *testing.T) {
	obj := Log{}
	obj.Level = FATAL
	obj.Message = "Dummy message"
	obj.Rq = &fasthttp.RequestCtx{}
	obj.Rq.Request.SetRequestURI("https://whatever/you/want?query=string&another=one")
	obj.Rq.Request.Header.Set("Referer", "https://www.google.co.uk")
	obj.Rq.Request.Header.Set("X-Forwarded-For", "8.8.8.8")
	output := createLog(obj)

	if output.ClientIP != "8.8.8.8" {
		t.Log(output.ClientIP)
		t.Error("Client IP is missing or wrong!")
	}

	if output.URL != "https://whatever/you/want" {
		t.Log(output.URL)
		t.Error("URL is missing or wrong!")
	}

	if output.QueryString != "query=string&another=one" {
		t.Log(output.QueryString)
		t.Error("QueryString is missing or wrong!")
	}

	if output.Referer != "https://www.google.co.uk" {
		t.Log(output.Referer)
		t.Error("Referer is missing or wrong!")
	}
}

func TestClientValuesAreEmpty(t *testing.T) {
	obj := Log{}
	obj.Level = FATAL
	obj.Message = "Dummy message"
	obj.Rq = nil
	output := createLog(obj)
	if output.ClientIP != "" || output.URL != "" || output.Referer != "" || output.QueryString != "" {
		t.Error("Client values should be empty!")
	}
}
