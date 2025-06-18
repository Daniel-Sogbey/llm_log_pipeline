package llm

import (
	"bytes"
	"encoding/json"
	"github.com/Daniel-Sogbey/llm_log_pipeline/internal/requester"
	"net/http"
)

type LLM struct {
	URL           string
	Authorization string
	Model         string
}

func (l *LLM) AnalyzeLog(request LLMRequestModel) (*LLMResponseModel, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := requester.MakeRequest[LLMResponseModel](http.MethodPost, l.URL, l.Authorization, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return response, err
}
