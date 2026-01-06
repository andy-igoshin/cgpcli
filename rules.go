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

func (cli *Cli) getComment(r []any) string {
	if len(r) > 4 {
		comment, _ := r[4].(string)
		return comment
	}
	return ""
}

func (cli *Cli) getMailRules(cmd, subject string) ([]MailRule, error) {
	raw, err := cli.getRawRules(cmd, subject)
	if err != nil {
		return nil, err
	}

	rules := make([]MailRule, 0, len(raw))
	for _, r := range raw {
		name, _ := r[1].(string)
		rule := MailRule{
			Priority:   cli.toInt(r[0]),
			Name:       name,
			Conditions: cli.toAnySlice(r[2]),
			Actions:    cli.toAnySlice(r[3]),
			Comment:    cli.getComment(r),
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (cli *Cli) getSignalRules(cmd, subject string) ([]SignalRule, error) {
	raw, err := cli.getRawRules(cmd, subject)
	if err != nil {
		return nil, err
	}

	rules := make([]SignalRule, 0, len(raw))
	for _, r := range raw {
		prio, delay, event := unpackSignalPriority(cli.toInt(r[0]))
		name, _ := r[1].(string)
		rule := SignalRule{
			Priority:   prio,
			Delay:      delay,
			EventCode:  event,
			Name:       name,
			Conditions: cli.toAnySlice(r[2]),
			Actions:    cli.toAnySlice(r[3]),
			Comment:    cli.getComment(r),
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// getRawRules — базовый метод для получения сырого списка правил от сервера.
func (cli *Cli) getRawRules(cmd, subject string) ([][]any, error) {
	if subject == "" {
		return nil, fmt.Errorf("parameter is required")
	}

	args := make([]any, 0, 2)
	args = append(args, cmd)
	if subject != "" {
		args = append(args, Atom(subject))
	}

	res, err := cli.Query(args...)
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

func setRulesGeneric[T ruleConstraint](cli *Cli, cmd, subject string, rules []T) error {
	if subject == "" || rules == nil {
		return fmt.Errorf("parameters are required")
	}

	rawRules := make([]any, len(rules))
	for i := range rules {
		rawRules[i] = rules[i].ToSlice()
	}

	args := make([]any, 0, 3)
	args = append(args, cmd)
	if subject != "" {
		args = append(args, Atom(subject))
	}
	args = append(args, rawRules)

	_, err := cli.Query(args...)
	return err
}

func updateRuleGeneric[T ruleConstraint](cli *Cli, cmd, subject string, rule T) error {
	if subject == "" {
		return fmt.Errorf("parameter is required")
	}
	if rule.GetName() == "" {
		return fmt.Errorf("rule name is required")
	}

	args := make([]any, 0, 3)
	args = append(args, cmd)
	if subject != "" {
		args = append(args, Atom(subject))
	}
	args = append(args, rule.ToSlice())

	_, err := cli.Query(args...)
	return err
}

// packSignalPriority упаковывает данные с учетом типа правила (служебное/пользовательское).
func packSignalPriority(name string, priority int, delay int, eventCode int) int {
	safePriority := priority % priorityScale

	// 1. Обработка служебных правил
	if strings.HasPrefix(name, "#") {
		// Разрешенные для Delay правила
		if name == "#ForkOnTimer" || name == "#MailOnTimer" {
			return (delayBase-delay)*priorityScale + safePriority
		}
		// Для остальных служебных правил используем delayBase (задержка 0)
		return delayBase*priorityScale + safePriority
	}

	// 2. Обработка пользовательских правил
	// Если указан eventCode (Ошибка/Занято/...), он имеет приоритет над Delay
	if eventCode == EventError || eventCode == EventBusy || eventCode == EventNoAnswer {
		return eventCode*priorityScale + safePriority
	}

	// Если события нет, пакуем обычную задержку
	if delay > 0 {
		return (delayBase-delay)*priorityScale + safePriority
	}

	return safePriority
}

// unpackSignalPriority извлекает приоритет, задержку или код события.
func unpackSignalPriority(packed int) (priority int, delay int, eventCode int) {
	priority = packed % priorityScale
	rawPrefix := packed / priorityScale

	// Если префикс 1, 3 или 5 — это коды событий (Error, Busy, NoAnswer)
	switch rawPrefix {
	case EventError, EventBusy, EventNoAnswer:
		eventCode = rawPrefix
		delay = 0
	default:
		// В остальных случаях это задержка (включая служебные правила)
		if rawPrefix > 0 {
			delay = delayBase - rawPrefix
		}
	}
	return
}
