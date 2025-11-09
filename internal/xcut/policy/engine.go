package policy

import (
	"context"
	"fmt"
)

type PolicyEngine struct {
	rules map[string]PolicyRule
}

type PolicyRule struct {
	Resource string
	Action   string
	Allow    func(subject, resource, action string) bool
}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{
		rules: make(map[string]PolicyRule),
	}
}

func (pe *PolicyEngine) AddRule(name string, rule PolicyRule) {
	pe.rules[name] = rule
}

func (pe *PolicyEngine) Evaluate(ctx context.Context, subject, resource, action string) (bool, error) {
	key := fmt.Sprintf("%s:%s", resource, action)
	if rule, exists := pe.rules[key]; exists {
		return rule.Allow(subject, resource, action), nil
	}
	return false, nil
}

func (pe *PolicyEngine) LoadDefaultRules() {
	pe.AddRule("country:read", PolicyRule{
		Resource: "country",
		Action:   "read",
		Allow:    func(subject, resource, action string) bool { return true },
	})
	pe.AddRule("country:write", PolicyRule{
		Resource: "country",
		Action:   "write",
		Allow:    func(subject, resource, action string) bool { return subject != "guest" },
	})
}