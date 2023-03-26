package entity

type Question struct {
	QuestionId  string `json:"questionId"`
	Description string `json:"description"`
	Title       string `json:"title"`
}

type QuestionPayload struct {
	QuestionId          string      `json:"questionId"`
	QuestionDescription string      `json:"questionDescription"`
	QuestionTitle       string      `json:"questionTitle"`
	CodeStub            *CodeStub   `json:"codeStub"`
	TestCase            []*TestCase `json:"testCases"`
}

type CodeStub struct {
	FunctionName   string    `json:"functionName"`
	Params         []*Params `json:"params"`
	ReturnDataType string    `json:"returnDataType"`
}

type Params struct {
	ParamName     string `json:"paramName"`
	ParamDataType string `json:"paramDataType"`
}
