package permission

type Mode string

const (
	ModeDefault Mode = "default"
	ModeAccept  Mode = "accept"
	ModePlan    Mode = "plan"
	ModeAuto    Mode = "auto"
)

type Behavior string

const (
	BehaviorAllow Behavior = "allow"
	BehaviorDeny  Behavior = "deny"
	BehaviorAsk   Behavior = "ask"
)

type Rule struct {
	ToolName string
	Pattern  string
	Behavior Behavior
}

type Decision struct {
	Behavior Behavior
	Reason   string
}

type Checker struct {
	mode  Mode
	rules []Rule
}

func NewChecker(mode Mode) *Checker {
	return &Checker{
		mode:  mode,
		rules: getDefaultRules(),
	}
}

func getDefaultRules() []Rule {
	return []Rule{
		{ToolName: "Read", Pattern: "*", Behavior: BehaviorAllow},
		{ToolName: "Glob", Pattern: "*", Behavior: BehaviorAllow},
		{ToolName: "Grep", Pattern: "*", Behavior: BehaviorAllow},
	}
}

func (c *Checker) SetRules(rules []Rule) {
	c.rules = rules
}

func (c *Checker) Check(toolName string, input map[string]interface{}) *Decision {
	if c.mode == ModeAccept {
		return &Decision{Behavior: BehaviorAllow, Reason: "accept mode"}
	}

	if c.mode == ModeAuto {
		return c.checkAutoMode(toolName, input)
	}

	for _, rule := range c.rules {
		if rule.ToolName != toolName {
			continue
		}

		if matchPattern(rule.Pattern, input) {
			return &Decision{Behavior: rule.Behavior, Reason: "rule matched"}
		}
	}

	if isDangerousCommand(toolName, input) {
		return &Decision{Behavior: BehaviorAsk, Reason: "dangerous command detected"}
	}

	return &Decision{Behavior: BehaviorAsk, Reason: "no matching rule"}
}

func (c *Checker) checkAutoMode(toolName string, input map[string]interface{}) *Decision {
	if isDangerousCommand(toolName, input) {
		return &Decision{Behavior: BehaviorAsk, Reason: "dangerous command in auto mode"}
	}

	return &Decision{Behavior: BehaviorAllow, Reason: "auto mode"}
}

func matchPattern(pattern string, input map[string]interface{}) bool {
	if pattern == "*" {
		return true
	}

	for _, value := range input {
		if strValue, ok := value.(string); ok {
			if matchStringPattern(pattern, strValue) {
				return true
			}
		}
	}

	return false
}

func matchStringPattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}

	return value == pattern || containsPattern(pattern, value)
}

func containsPattern(pattern, value string) bool {
	return len(pattern) > 0 && len(value) >= len(pattern) && value[:len(pattern)] == pattern
}
