package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"git.vsu.ru/ai/cgpcli"
)

// Подготовка данных
var (
	benchData    = generateLargeMap(5000) // 5000 ключей, примерно 10-15к элементов
	cgpRawBytes  []byte
	jsonRawBytes []byte
)

func generateLargeMap(size int) map[string]any {
	m := make(map[string]any, size)
	for i := 0; i < size; i++ {
		m[fmt.Sprintf("key_header_%d", i)] = "some_standard_string_value"
		m[fmt.Sprintf("key_list_%d", i)] = []any{"Read", i, true, "user@domain.tld"}
		m[fmt.Sprintf("key_sub_%d", i)] = map[string]any{
			"enabled": true,
			"count":   int64(i),
		}
	}
	return m
}

func init() {
	s, _ := cgpcli.MarshalToString(benchData)
	cgpRawBytes = []byte(s)
	jsonRawBytes, _ = json.Marshal(benchData)
}

// --- CGP BENCHMARKS ---

func BenchmarkMarshalCGP_ToString(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = cgpcli.MarshalToString(benchData)
	}
}

// Прямой маршалинг в Writer (аналог работы вашего Cli.Send)
func BenchmarkMarshalCGP_ToWriter(b *testing.B) {
	b.ReportAllocs()
	// Используем discard, чтобы измерить только скорость маршалинга без I/O
	bw := bufio.NewWriter(io.Discard)
	for i := 0; i < b.N; i++ {
		_ = cgpcli.MarshalTo(bw, benchData)
		bw.Flush()
		bw.Reset(io.Discard)
	}
}

func BenchmarkUnmarshalCGP(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var res any
		_ = cgpcli.Unmarshal(cgpRawBytes, &res)
	}
}

// --- JSON BENCHMARKS ---

func BenchmarkMarshalJSON_Standard(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(benchData)
	}
}

func BenchmarkUnmarshalJSON_Standard(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var res any
		_ = json.Unmarshal(jsonRawBytes, &res)
	}
}
