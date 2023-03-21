package test_use_case

import (
	"server/entity"
)

type TestUseCase interface {
	CreateNewTestCase(model *entity.TestCase) error
}
