package cgpcli

import (
	"fmt"
	"strings"
)

const (
	delayBase     = 1000
	priorityScale = 100

	// Константы событий для пользовательских правил
	EventError    = 1 // 100
	EventBusy     = 3 // 300
	EventNoAnswer = 5 // 500
)

// -----------------------------------------------------------------------------
// Types & Interfaces
// -----------------------------------------------------------------------------

type MailRule struct {
	Priority   int    `json:"priority"`
	Name       string `json:"name"`
	Conditions []any  `json:"conditions"`
	Actions    []any  `json:"actions"`
	Comment    string `json:"comment,omitempty"`
}

type SignalRule struct {
	Priority   int    `json:"priority"`
	Delay      int    `json:"delay"`
	EventCode  int    `json:"event_code,omitempty"`
	Name       string `json:"name"`
	Conditions []any  `json:"conditions"`
	Actions    []any  `json:"actions"`
	Comment    string `json:"comment,omitempty"`
}

type ruleConstraint interface {
	MailRule | SignalRule
	ToSlice() []any
	GetName() string
}

func (r MailRule) GetName() string   { return r.Name }
func (r SignalRule) GetName() string { return r.Name }

func (r MailRule) ToSlice() []any {
	slice := []any{r.Priority, r.Name, r.Conditions, r.Actions}
	if r.Comment != "" {
		slice = append(slice, r.Comment)
	}
	return slice
}

func (r SignalRule) ToSlice() []any {
	packed := packSignalPriority(r.Name, r.Priority, r.Delay, r.EventCode)
	slice := []any{packed, r.Name, r.Conditions, r.Actions}
	if r.Comment != "" {
		slice = append(slice, r.Comment)
	}
	return slice
}

// -----------------------------------------------------------------------------
// Internal Helpers (Private)
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
		_, err = cli.Query(cmd, Atom(subject), rawRules)
	} else {
		_, err = cli.Query(cmd, rawRules)
	}
	return err
}

func updateRuleGeneric[T ruleConstraint](cli *Cli, cmd, subject string, rule T) error {
	if rule.GetName() == "" {
		return fmt.Errorf("%s: rule name is required", cmd)
	}

	var err error
	if subject != "" {
		_, err = cli.Query(cmd, Atom(subject), rule.ToSlice())
	} else {
		_, err = cli.Query(cmd, rule.ToSlice())
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
