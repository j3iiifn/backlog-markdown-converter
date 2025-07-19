package main

import (
	"testing"
)

func TestVersionCommand(t *testing.T) {
	// バージョン変数の値をテスト
	expected := "0.1.0"
	actual := version

	if actual != expected {
		t.Errorf("Expected %q, but got %q", expected, actual)
	}
}
