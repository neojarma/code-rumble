package test_repository

import "executor/entity"

type TestResultRepository interface {
	CreateTestResult(results []*entity.TestResult) error
}
