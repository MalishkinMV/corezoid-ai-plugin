package main

import (
	"io/ioutil"
	"testing"
)

func TestManualValidation_ApiCopyToItself(t *testing.T) {
	// Читаем тестовый JSON
	data, err := ioutil.ReadFile("./samples/api_copy_to_itself.json")
	if err != nil {
		t.Fatalf("failed to read sample file: %v", err)
	}

	// Создаем Executor с нужным ProcessID (например, 123)
	executor := &Executor{ProcessID: 123}

	// Проверяем, что BeforeValidation возвращает ошибку
	err = executor.BeforeValidation(string(data), nil)
	if err == nil {
		t.Errorf("expected error for api_copy to itself, got nil")
	}
	t.Logf("Expected error: %v", err)
}
