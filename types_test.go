package cgpcli

import (
	"bufio"
	"bytes"
	"net"
	"testing"
)

func TestCGPIP_WriteCGP(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		port int
		want string
	}{
		{"IPv4_NoPort", net.ParseIP("127.0.0.1"), 0, "#I[127.0.0.1]"},
		{"IPv4_WithPort", net.ParseIP("10.0.0.1"), 25, "#I[10.0.0.1]:25"},
		{"IPv6_NoPort", net.ParseIP("2001:db8::1"), 0, "#I[2001:db8::1]"},
		{"IPv6_WithPort", net.ParseIP("::1"), 8080, "#I[::1]:8080"},
		{"NilIP", nil, 0, "#NULL#"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CGPIP{IP: tt.ip, Port: tt.port}

			var b bytes.Buffer
			// Используем bufio.Writer, так как он удовлетворяет интерфейсу encodeWriter
			w := bufio.NewWriter(&b)

			err := c.WriteCGP(w)
			w.Flush() // Важно сбросить буфер в bytes.Buffer

			if err != nil {
				t.Fatalf("WriteCGP() error = %v", err)
			}

			if got := b.String(); got != tt.want {
				t.Errorf("WriteCGP() output = %v, want %v", got, tt.want)
			}
		})
	}
}
