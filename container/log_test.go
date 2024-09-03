package main

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
)

// Testa a geração de mensagens de log informativas comparando a saída gerada
// com a saída esperada
func TestLogInfo(t *testing.T) {
	input := "a b c"
	correlation := "123456789"
	want := regexp.MustCompile(fmt.Sprintf(`\[INFO\] \{%s\} %s\n`, correlation, input))
	buffer := &bytes.Buffer{}
	log := NewLog(buffer)
	log.SetCorrelation(correlation)
	log.Info(input)
	output := buffer.String()
	if !want.MatchString(output) {
		t.Fatalf("output invalid, received {%s}", output)
	}
}

// Testa a geração de mensagens de log de alerta comparando a saída gerada
// com a saída esperada
func TestLogWarning(t *testing.T) {
	input := "a b c"
	correlation := "123456789"
	want := regexp.MustCompile(fmt.Sprintf(`\[WARN\] \{%s\} %s\n`, correlation, input))
	buffer := &bytes.Buffer{}
	log := NewLog(buffer)
	log.SetCorrelation(correlation)
	log.Warning(input)
	output := buffer.String()
	if !want.MatchString(output) {
		t.Fatalf("output invalid, received {%s}", output)
	}
}

// Testa a geração de mensagens de log de erro comparando a saída gerada
// com a saída esperada
func TestLogError(t *testing.T) {
	input := "a b c"
	correlation := "123456789"
	want := regexp.MustCompile(fmt.Sprintf(`\[ERROR\] \{%s\} %s\n`, correlation, input))
	buffer := &bytes.Buffer{}
	log := NewLog(buffer)
	log.SetCorrelation(correlation)
	log.Error(input)
	output := buffer.String()
	if !want.MatchString(output) {
		t.Fatalf("output invalid, received {%s}", output)
	}
}
