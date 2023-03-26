package test_cases_repository

import "server/entity"

type TestCaseRepository interface {
	CreateNewTestCase(model *entity.TestCase) error
	CreateBulkTestCase(req []*entity.TestCase) error
}
