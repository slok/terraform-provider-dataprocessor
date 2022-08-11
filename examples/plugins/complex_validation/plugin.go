package tf

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type Options struct {
	CheckRunbook  bool
	CheckSeverity bool
	CheckTeam     bool
}

type PromRulesRoot struct {
	Groups []PromRuleGroup `json:"groups"`
}
type PromRuleGroup struct {
	Name  string     `json:"name"`
	Rules []PromRule `json:"rules"`
}

type PromRule struct {
	Name        string            `json:"name"`  // Recording rules.
	Alert       string            `json:"alert"` // Alerting rules.
	Annotations map[string]string `json:"annotations"`
	Labels      map[string]string `json:"labels"`
}

type RuleGroupValidator struct {
	validator Validator
}

// ProcessorPluginV1 will check that all prometheus rules meet the minimum requirements.
func ProcessorPluginV1(ctx context.Context, inputData string, vars map[string]string) (string, error) {
	// Load options
	opts, err := loadOptions(vars)
	if err != nil {
		return "", fmt.Errorf("could not load options: %w", err)
	}

	// Load rules.
	root := []PromRulesRoot{}
	err = json.Unmarshal([]byte(inputData), &root)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal role groups: %w", err)
	}

	// Validate.
	validator := NewRuleGroupValidator(*opts)
	invalidMsgs := map[string][]string{}
	invalid := false
	for _, rg := range root {
		// Validate every group.
		for _, rg := range rg.Groups {
			msgs, err := validator.ValidateRuleGroup(ctx, rg)
			if err != nil {
				return "", fmt.Errorf("could not validate rule group %q: %w", rg.Name, err)
			}

			if len(msgs) > 0 {
				invalid = true
			}

			invalidMsgs[rg.Name] = msgs
		}
	}

	if !invalid {
		return "valid", nil
	}

	// Create the invalid message in pretty format.
	msg := "\n"
	for k, v := range invalidMsgs {
		if len(v) == 0 {
			msg += fmt.Sprintf("✔️ %s:\n", k)
			continue
		}

		msg += fmt.Sprintf("❌ %s:\n", k)
		for _, v := range v {
			msg += fmt.Sprintf("  ⭕ %s\n", v)
		}
	}

	return msg, fmt.Errorf(msg)
}

func loadOptions(vars map[string]string) (*Options, error) {
	var opts Options
	var err error

	val, ok := vars["check_runbook"]
	if ok {
		opts.CheckRunbook, err = strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
	}

	val, ok = vars["check_severity"]
	if ok {
		opts.CheckSeverity, err = strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
	}

	val, ok = vars["check_team"]
	if ok {
		opts.CheckTeam, err = strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
	}

	return &opts, nil
}

func NewRuleGroupValidator(options Options) *RuleGroupValidator {
	valChain := []Validator{}
	if options.CheckRunbook {
		valChain = append(valChain, NewValidAlertRunbook())
	}

	if options.CheckSeverity {
		valChain = append(valChain, NewValidAlertSeverity())
	}

	if options.CheckTeam {
		valChain = append(valChain, NewValidAlertTeam())
	}

	return &RuleGroupValidator{
		validator: NewValidatorChain(valChain...),
	}
}

func (r RuleGroupValidator) ValidateRuleGroup(ctx context.Context, rg PromRuleGroup) (invalidMsg []string, err error) {
	msgs := []string{}
	for _, pr := range rg.Rules {
		name := pr.Name
		if name == "" {
			name = pr.Alert
		}

		warns, err := r.validatePromRule(ctx, pr)
		if err != nil {
			return nil, fmt.Errorf("could not validate %q prom rule: %w", name, err)
		}

		for _, w := range warns {
			msgs = append(msgs, name+": "+w)
		}
	}

	return msgs, nil
}

func (r RuleGroupValidator) validatePromRule(ctx context.Context, pr PromRule) (invalidMsg []string, err error) {
	msgs := []string{}

	msg, valid, err := r.validator.ValidatePromRule(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("could not validate rule: %w", err)
	}

	if !valid {
		msgs = append(msgs, msg...)
	}

	return msgs, nil
}

// Validator knows how to validate a prometheus rules, it will return info messages to explain
// why the rule is not valid.
type Validator interface {
	ValidatePromRule(ctx context.Context, pr PromRule) (msgs []string, valid bool, err error)
}

// ValidatorFunc is a helper type to create validators without the need to declare a new type.
type ValidatorFunc func(ctx context.Context, pr PromRule) (msgs []string, valid bool, err error)

func (v ValidatorFunc) ValidatePromRule(ctx context.Context, pr PromRule) (msg []string, valid bool, err error) {
	return v(ctx, pr)
}

// NewValidatorChain is a validator that will execute a list of validator as a chain, this way
// with this simple Validator we can execute a group of validators one after another.
func NewValidatorChain(vs ...Validator) Validator {
	return ValidatorFunc(func(ctx context.Context, pr PromRule) (msg []string, valid bool, err error) {
		allMsgs := []string{}

		isValid := true
		for _, v := range vs {
			msg, valid, err := v.ValidatePromRule(ctx, pr)
			if err != nil {
				return nil, false, err
			}

			if !valid {
				isValid = false
				allMsgs = append(allMsgs, msg...)
			}
		}

		return allMsgs, isValid, nil
	})
}

// NewValidAlertRunbook returns a validator that will check that the alert has a runbook and
// is a proper URL.
func NewValidAlertRunbook() Validator {
	return ValidatorFunc(func(ctx context.Context, pr PromRule) (msg []string, valid bool, err error) {
		if pr.Alert == "" {
			return nil, true, nil
		}

		runbookURL, ok := pr.Annotations["runbook_url"]
		if !ok {
			return []string{"Alert runbook missing"}, false, nil
		}

		_, err = url.ParseRequestURI(runbookURL)
		if err != nil {
			return []string{fmt.Sprintf("Alert runbook %q URL is not valid: %s", runbookURL, err)}, false, nil
		}

		return nil, true, nil
	})
}

var validSeverities = map[string]struct{}{
	"info":     {},
	"warning":  {},
	"critical": {},
	"page":     {},
}

// NewValidAlertSeverity returns a validator that will check that the alert has a runbook and
// is a proper URL.
func NewValidAlertSeverity() Validator {
	return ValidatorFunc(func(ctx context.Context, pr PromRule) (msgs []string, valid bool, err error) {
		if pr.Alert == "" {
			return nil, true, nil
		}

		severity, ok := pr.Labels["severity"]
		if !ok {
			return []string{"Alert severity missing"}, false, nil
		}

		_, ok = validSeverities[severity]
		if !ok {
			return []string{fmt.Sprintf("Alert severity %q is not supported", severity)}, false, nil
		}

		return nil, true, nil
	})
}

// NewValidAlertTeam returns a validator that checks the label team is present on the alerts.
func NewValidAlertTeam() Validator {
	return ValidatorFunc(func(ctx context.Context, pr PromRule) (msg []string, valid bool, err error) {
		if pr.Alert == "" {
			return nil, true, nil
		}

		_, ok := pr.Labels["team"]
		if !ok {
			return []string{"The team is missing"}, false, nil
		}

		return nil, true, nil
	})
}
