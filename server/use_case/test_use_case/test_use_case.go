package test_use_case

import (
	"server/entity"
)

type TestUseCase interface {
	CreateNewTestCase(model *entity.TestCase) error
	CreateBulkTestCase(req []*entity.TestCase, questionId string) error
}
