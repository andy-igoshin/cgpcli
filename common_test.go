package cgpcli

import (
	"os"
	"strings"
	"testing"
)

func TestNew_host_fail(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	_, err := New("127.0.0.2", cgpuser, cgppass, Plain, false)
	if err != nil && !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("%s", err)
	}
}

func TestNew_login_fail(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	_, err := New("127.0.0.1", cgpuser, cgppass+"1", Plain, false)
	if err != nil && !strings.Contains(err.Error(), "incorrect password or account name") {
		t.Errorf("%s", err)
	}
}

func TestNew_login_plain(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, Plain, false)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Logout()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestNew_login_apop(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Logout()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestNew_login_stls(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, Plain, true)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Logout()
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestNew_logout(t *testing.T) {
	cgpuser := os.Getenv("CGPUSER")
	cgppass := os.Getenv("CGPPASS")
	if cgpuser == "" || cgppass == "" {
		t.Errorf("specify environment variables CGPUSER and CGPPPASS")
	}

	cli, err := New("127.0.0.1", cgpuser, cgppass, APOP, false)
	if err != nil {
		t.Errorf("%s", err)
	}

	err = cli.Logout()
	if err != nil {
		t.Errorf("%s", err)
	}
}
