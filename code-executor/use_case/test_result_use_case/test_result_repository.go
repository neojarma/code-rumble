package test_use_case

import "executor/entity"

type TestResultUseCase interface {
	CreateTestResult(results []*entity.TestResult) error
}
