package main

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Define a estrura do log
type Log struct {
	mu          sync.Mutex
	o           io.Writer
	correlation string
}

// Cria o log com a saída desejada, se for vazio, então usa a saída padrão
func NewLog(o io.Writer) *Log {
	if o == nil {
		o = os.Stdout
	}
	return &Log{o: o}
}

// Escreve na saida a mensagem
func (p *Log) log(message string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	b := []byte(fmt.Sprintf("%s\n", message))
	_, err := p.o.Write(b)
	if err != nil {
		if p.o != os.Stdout {
			os.Stdout.Write([]byte(fmt.Sprintf("failed to write log, %s", err)))
			os.Stdout.Write(b)
		}
	}
}

// Escreve mensagem padrão
func (p *Log) Print(format string, a ...any) {
	p.log(fmt.Sprintf(format, a...))
}

// Escreve mensagem informativa
func (p *Log) Info(format string, a ...any) {
	p.log(fmt.Sprintf("[INFO] %s%s", p.printCorrelation(), fmt.Sprintf(format, a...)))
}

// Escreve mensagem aviso
func (p *Log) Warning(format string, a ...any) {
	p.log(fmt.Sprintf("[WARN] %s%s", p.printCorrelation(), fmt.Sprintf(format, a...)))
}

// Escreve mensagem erro
func (p *Log) Error(format string, a ...any) {
	p.log(fmt.Sprintf("[ERROR] %s%s", p.printCorrelation(), fmt.Sprintf(format, a...)))
}

// Escreve mensagem erro e encerra a execução
func (p *Log) Fatal(format string, a ...any) {
	p.log(fmt.Sprintf("[ERROR] %s%s", p.printCorrelation(), fmt.Sprintf(format, a...)))
	os.Exit(1)
}

// Define o Correlation para as mensagens de logs que serão geradas
func (p *Log) SetCorrelation(correlation string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.correlation = correlation
}

// Retorna o Correlation atual
func (p *Log) Correlation() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.correlation
}

// Retorna o Correlation formatado para impressão
func (p *Log) printCorrelation() string {
	if p.correlation != "" {
		return fmt.Sprintf("{%s} ", p.correlation)
	}
	return ""
}
