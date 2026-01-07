package cgpcli

import (
	"reflect"
	"testing"
)

func TestSignalPriorityArithmetic(t *testing.T) {
	tests := []struct {
		name      string
		prio      int
		delay     int
		event     int
		ruleName  string
		want      int
	}{
		{"Delay 0 (Base case)", 5, 0, 0, "#Block", 100005},
		{"Delay 10 minutes", 1, 10, 0, "#ForkOnTimer", 99001},
		{"Event Code Busy", 5, 0, EventBusy, "UserRule", 305},
		{"Event Code Error", 12, 0, EventError, "UserRule", 112},
		{"Normal Priority (No prefix)", 50, 0, 0, "NormalRule", 50},
		{"Max Priority wrapping", 150, 0, 0, "NormalRule", 50}, // 150 % 100
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := packSignalPriority(tt.ruleName, tt.prio, tt.delay, tt.event)
			if got != tt.want {
				t.Errorf("pack() = %v, want %v", got, tt.want)
			}

			// Проверка обратимости
			p, d, e := unpackSignalPriority(got)
			if p != (tt.prio % priorityScale) {
				t.Errorf("unpack() priority = %v, want %v", p, tt.prio%priorityScale)
			}
			// Для Event префикс вытесняет Delay в логике unpack
			if tt.event > 0 {
				if e != tt.event {
					t.Errorf("unpack() event = %v, want %v", e, tt.event)
				}
			} else if d != tt.delay {
				t.Errorf("unpack() delay = %v, want %v", d, tt.delay)
			}
		})
	}
}

func TestRulesToSlice(t *testing.T) {
	t.Run("MailRule with Comment", func(t *testing.T) {
		r := MailRule{
			Priority: 5, Name: "M1", 
			Conditions: []any{"Subject", "is", "test"}, 
			Actions: []any{"Discard"}, 
			Comment: "Hide it",
		}
		got := r.ToSlice()
		want := []any{5, "M1", r.Conditions, r.Actions, "Hide it"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("SignalRule Packing", func(t *testing.T) {
		r := SignalRule{
			Priority: 5, Delay: 0, Name: "#Block",
			Conditions: []any{"From", "is", "spam"},
			Actions: []any{"Reject"},
		}
		got := r.ToSlice()
		// 1000*100 + 5 = 100005
		if got[0] != 100005 {
			t.Errorf("Expected packed priority 100005, got %v", got[0])
		}
	})
}
