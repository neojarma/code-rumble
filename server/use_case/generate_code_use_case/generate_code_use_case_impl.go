package generate_code

import "server/entity"

type GenerateCode interface {
	GenerateJavaScriptCode(req *entity.CodeStub) ([]byte, []byte, error)
}
