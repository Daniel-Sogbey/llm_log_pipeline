package llm

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Daniel-Sogbey/llm_log_pipeline/internal/requester"
)

type TogetherAi struct {
	URL           string
	Authorization string
	Model         string
}

func (t *TogetherAi) AnalyzeLog(request TogetherAiLLMRequestModel) (*TogetherAiLLMResponseModel, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := requester.MakeRequest[TogetherAiLLMResponseModel](http.MethodPost, t.URL, t.Authorization,
		bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return response, err
}
