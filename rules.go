// # Rules Administration
//
// The Rules section provides specialized management for message and signal
// processing logic. CommuniGate Pro uses a rule-based engine to automate
// actions like filtering, forwarding, or complex VoIP call routing.
//
// Key capabilities include:
//   - Dual-Engine Support: separate structures and methods for [MailRule]
//     (standard mail processing) and [SignalRule] (SIP/Signal processing).
//   - Priority Logic: automatic handling of "packed" priorities in Signal
//     rules, where the library transparently manages delays and event codes
//     (e.g., Error, Busy, No Answer).
//   - Flexible Conditions & Actions: rules are represented as structured
//     slices, allowing for complex nested logic as defined in CommuniGate Pro.
//   - Generics-Powered Operations: unified internal logic for setting and
//     updating rules using Go Generics to ensure type safety across
//     different rule levels.
package cgpcli

import (
	"fmt"
	"strings"
)

const (
	delayBase     = 1000
	priorityScale = 100

	EventError    = 1 // 100
	EventBusy     = 3 // 300
	EventNoAnswer = 5 // 500
)

// MailRule represents a standard email processing rule.
type MailRule struct {
	Priority   int    // execution order (0 to 99).
	Name       string // unique rule identifier.
	Conditions []any  // filter criteria (e.g., ["From", "is", "spammer@..."]).
	Actions    []any  // resulting actions (e.g., ["Discard"]).
	Comment    string // optional administrative note.
}

// SignalRule represents a Real-Time Application (SIP/Signal) rule.
type SignalRule struct {
	Priority   int    // base priority (0 to 99).
	Delay      int    // delay in seconds before the rule is applied.
	EventCode  int    // specific event trigger (EventError, EventBusy, EventNoAnswer).
	Name       string // unique rule identifier.
	Conditions []any  // signal-specific conditions.
	Actions    []any  // signaling actions (e.g., ["Redirect", "target"]).
	Comment    string // optional administrative note.
}

type ruleConstraint interface {
	MailRule | SignalRule
	ToSlice() []any
	GetName() string
}

// GetName returns the unique identifier of the rule.
//
// Returns:
//   - string: the name of the [MailRule].
func (r MailRule) GetName() string {
	return r.Name
}

// GetName returns the unique identifier of the rule.
//
// Returns:
//   - string: the name of the [SignalRule].
func (r SignalRule) GetName() string {
	return r.Name
}

// ToSlice converts the rule structure into a raw slice format compatible
// with the CommuniGate Pro CLI protocol.
//
// Returns:
//   - []any: a slice containing [Priority, Name, Conditions, Actions, (Comment)].
func (r MailRule) ToSlice() []any {
	slice := []any{r.Priority, r.Name, r.Conditions, r.Actions}
	if r.Comment != "" {
		slice = append(slice, r.Comment)
	}
	return slice
}

// ToSlice converts the rule structure into a raw slice format compatible
// with the CommuniGate Pro CLI protocol.
//
// For [SignalRules], this method automatically performs priority packing
// (combining priority, delay, and event code).
//
// Returns:
//   - []any: a slice containing [Priority, Name, Conditions, Actions, (Comment)].
func (r SignalRule) ToSlice() []any {
	packed := packSignalPriority(r.Name, r.Priority, r.Delay, r.EventCode)
	slice := []any{packed, r.Name, r.Conditions, r.Actions}
	if r.Comment != "" {
		slice = append(slice, r.Comment)
	}
	return slice
}

// -----------------------------------------------------------------------------
// Internal Helpers
// -----------------------------------------------------------------------------

func (cli *Cli) getMailRules(cmd, subject string) ([]MailRule, error) {
	raw, err := cli.getRawRules(cmd, subject)
	if err != nil {
		return nil, err
	}

	rules := make([]MailRule, 0, len(raw))
	for _, r := range raw {
		name, _ := r[1].(string)
		rules = append(rules, MailRule{
			Priority:   toInt(r[0]),
			Name:       name,
			Conditions: toAnySlice(r[2]),
			Actions:    toAnySlice(r[3]),
			Comment:    getRuleComment(r),
		})
	}
	return rules, nil
}

func (cli *Cli) getRawRules(cmd, subject string) ([][]any, error) {
	var res any
	var err error

	if subject != "" {
		res, err = cli.Query(cmd, Atom(subject))
	} else {
		res, err = cli.Query(cmd)
	}

	if err != nil {
		return nil, err
	}

	rawList, ok := res.([]any)
	if !ok {
		return [][]any{}, nil
	}

	result := make([][]any, 0, len(rawList))
	for _, item := range rawList {
		if r, ok := item.([]any); ok && len(r) >= 4 {
			result = append(result, r)
		}
	}
	return result, nil
}

func (cli *Cli) getSignalRules(cmd, subject string) ([]SignalRule, error) {
	raw, err := cli.getRawRules(cmd, subject)
	if err != nil {
		return nil, err
	}

	rules := make([]SignalRule, 0, len(raw))
	for _, r := range raw {
		prio, delay, event := unpackSignalPriority(toInt(r[0]))
		name, _ := r[1].(string)
		rules = append(rules, SignalRule{
			Priority:   prio,
			Delay:      delay,
			EventCode:  event,
			Name:       name,
			Conditions: toAnySlice(r[2]),
			Actions:    toAnySlice(r[3]),
			Comment:    getRuleComment(r),
		})
	}
	return rules, nil
}

func getRuleComment(r []any) string {
	if len(r) > 4 {
		comment, _ := r[4].(string)
		return comment
	}
	return ""
}

// -----------------------------------------------------------------------------
// Generics Logic
// -----------------------------------------------------------------------------

func setRulesGeneric[T ruleConstraint](cli *Cli, cmd, subject string, rules []T) error {
	if rules == nil {
		return fmt.Errorf("%s: rules list is required", cmd)
	}

	rawRules := make([]any, len(rules))
	for i := range rules {
		rawRules[i] = rules[i].ToSlice()
	}

	var err error
	if subject != "" {
		err = cli.QueryNV(cmd, Atom(subject), rawRules)
	} else {
		err = cli.QueryNV(cmd, rawRules)
	}
	return err
}

func updateRuleGeneric[T ruleConstraint](cli *Cli, cmd, subject string, rule T) error {
	if rule.GetName() == "" {
		return fmt.Errorf("%s: rule name is required", cmd)
	}

	var err error
	if subject != "" {
		err = cli.QueryNV(cmd, Atom(subject), rule.ToSlice())
	} else {
		err = cli.QueryNV(cmd, rule.ToSlice())
	}
	return err
}

// -----------------------------------------------------------------------------
// Signal Priority Packing Logic
// -----------------------------------------------------------------------------

func packSignalPriority(name string, priority int, delay int, eventCode int) int {
	safePriority := priority % priorityScale

	if strings.HasPrefix(name, "#") {
		if name == "#ForkOnTimer" || name == "#MailOnTimer" {
			return (delayBase-delay)*priorityScale + safePriority
		}
		return delayBase*priorityScale + safePriority
	}

	if eventCode == EventError || eventCode == EventBusy || eventCode == EventNoAnswer {
		return eventCode*priorityScale + safePriority
	}

	if delay > 0 {
		return (delayBase-delay)*priorityScale + safePriority
	}

	return safePriority
}

func unpackSignalPriority(packed int) (priority int, delay int, eventCode int) {
	priority = packed % priorityScale
	rawPrefix := packed / priorityScale

	switch rawPrefix {
	case EventError, EventBusy, EventNoAnswer:
		eventCode = rawPrefix
		delay = 0
	default:
		if rawPrefix > 0 {
			delay = delayBase - rawPrefix
		}
	}
	return
}
