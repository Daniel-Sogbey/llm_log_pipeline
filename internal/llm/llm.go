package llm

type LLM[I any, O any] interface {
	AnalyzeLog(request I) (response O, err error)
}
