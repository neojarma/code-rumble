package generate_code

import (
	"bytes"
	"fmt"
	"os"
	"server/entity"
)

type GenerateCodeImpl struct {
}

func NewGenerateCode() GenerateCode {
	return &GenerateCodeImpl{}
}

func (g *GenerateCodeImpl) GenerateJavaScriptCode(req *entity.CodeStub) ([]byte, []byte, error) {
	runCode, err := generateRunCodeForJs(req)
	if err != nil {
		return nil, nil, err
	}

	stubCode, err := generateStubCode(req)
	if err != nil {
		return nil, nil, err
	}

	return runCode, stubCode, nil
}

func generateStubCode(req *entity.CodeStub) ([]byte, error) {
	b, err := readMasterFile("./master/js_master_answer.txt")
	if err != nil {
		return nil, err
	}

	replacedFunction := bytes.ReplaceAll(b, []byte("{{f_name}}"), []byte(req.FunctionName))

	paramStr := ""
	for i, v := range req.Params {
		if i == len(req.Params)-1 {
			paramStr += v.ParamName
		} else {
			paramStr += fmt.Sprintf("%s,", v.ParamName)
		}
	}

	replacedParams := bytes.ReplaceAll(replacedFunction, []byte("{{params}}"), []byte(paramStr))

	return replacedParams, nil
}

func generateRunCodeForJs(req *entity.CodeStub) ([]byte, error) {
	b, err := readMasterFile("./master/js_master_question.txt")
	if err != nil {
		return nil, err
	}

	returnStr := determineJsDataType(req.ReturnDataType, "test[`output`]")

	paramsStr := ""
	for i, v := range req.Params {
		datatype := determineJsDataType(v.ParamDataType, fmt.Sprintf("rawInput[%v]", i))
		paramsStr += fmt.Sprintf("%s,", datatype)
	}

	afterReplaceReturn := bytes.ReplaceAll(b, []byte("{{parse_return_replace}}"), []byte(returnStr))
	afterReplaceFunctionName := bytes.ReplaceAll(afterReplaceReturn, []byte("{{function_name_replace}}"), []byte(req.FunctionName))
	replaceParam := bytes.ReplaceAll(afterReplaceFunctionName, []byte("{{loop_params_replace}}"), []byte(paramsStr))

	return replaceParam, nil
}

func readMasterFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func determineJsDataType(str, replacer string) string {
	switch str {
	case "float":
		return fmt.Sprintf("parseFloat(%s)", replacer)
	case "int":
		return fmt.Sprintf("parseInt(%s)", replacer)
	case "bool":
		return fmt.Sprintf("(%s === true) ? true : false", replacer)
	case "string":
		return replacer
	case "[]float":
		return fmt.Sprintf("%s.split(` `).map(v => parseFloat(v))", replacer)
	case "[]int":
		return fmt.Sprintf("%s.split(` `).map(v => parseInt(v))", replacer)
	case "[]str":
		return fmt.Sprintf("%s.split(` `)", replacer)
	case "[]bool":
		return fmt.Sprintf("%s.split(` `).map(v => (v === `true`) ? true : false)", replacer)
	default:
		return ""
	}
}
