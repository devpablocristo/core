package risk

import (
	"slices"
	"strings"
	"time"

	kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"
)

const (
	thresholdAllow       = 0.5
	thresholdEnhancedLog = 1.0
)

type Config struct {
	Thresholds            Thresholds
	HighActions           []string
	MediumActions         []string
	BusinessHours         BusinessHours
	FrequencyThresholds   FrequencyThresholds
	ActorThresholds       ActorThresholds
	SuccessRateThresholds SuccessRateThresholds
	Amplifications        []Amplification
	SensitiveSystems      []string
}

type Thresholds struct {
	Allow            float64
	EnhancedLog      float64
	RequireApproval  float64
	Deny             float64
	MaxAmplification float64
}

type BusinessHours struct {
	Start int
	End   int
}

type FrequencyThresholds struct {
	Warning  int
	Critical int
}

type ActorThresholds struct {
	Unknown int
	New     int
}

type SuccessRateThresholds struct {
	Low       float64
	Moderate  float64
	Excellent float64
}

type Amplification struct {
	Factors    []string
	Multiplier float64
	Reason     string
}

type History struct {
	ActorHistory    int
	RecentFrequency int
	SuccessRate     float64
}

func DefaultConfig() Config {
	return Config{
		Thresholds: Thresholds{
			Allow:            0.5,
			EnhancedLog:      1.0,
			RequireApproval:  1.5,
			Deny:             2.0,
			MaxAmplification: 3.0,
		},
		HighActions:   []string{"alert.silence", "runbook.execute", "delete"},
		MediumActions: []string{"incident.resolve", "config.update", "deploy.trigger"},
		BusinessHours: BusinessHours{Start: 9, End: 18},
		FrequencyThresholds: FrequencyThresholds{
			Warning:  10,
			Critical: 20,
		},
		ActorThresholds: ActorThresholds{
			Unknown: 0,
			New:     10,
		},
		SuccessRateThresholds: SuccessRateThresholds{
			Low:       0.5,
			Moderate:  0.8,
			Excellent: 0.95,
		},
		Amplifications: []Amplification{
			{Factors: []string{"off_hours", "actor_unknown"}, Multiplier: 1.8, Reason: "off-hours + unknown actor"},
			{Factors: []string{"action_type", "frequency_anomaly"}, Multiplier: 1.5, Reason: "risky action + frequency anomaly"},
			{Factors: []string{"actor_unknown", "target_sensitivity"}, Multiplier: 1.6, Reason: "unknown actor + sensitive target"},
			{Factors: []string{"off_hours", "actor_unknown", "frequency_anomaly"}, Multiplier: 2.5, Reason: "full cascade: off-hours + unknown + frequency"},
			{Factors: []string{"action_type", "off_hours", "target_sensitivity"}, Multiplier: 2.0, Reason: "risky action + off-hours + sensitive target"},
		},
		SensitiveSystems: []string{"production", "prod"},
	}
}

func Evaluate(request kerneldomain.Request, history History, config Config, policyRiskOverride *kerneldomain.RiskLevel, now time.Time) kerneldomain.RiskAssessment {
	config = normalizeConfig(config)
	if now.IsZero() {
		now = time.Now().UTC()
	}

	factors := evaluateFactors(request, history, config, now)
	rawScore := sumFactors(factors)
	amplification := calculateAmplification(factors, config)
	finalScore := rawScore * amplification
	if policyRiskOverride != nil {
		finalScore = applyPolicyOverride(*policyRiskOverride, finalScore, config)
	}

	return kerneldomain.RiskAssessment{
		Factors:       factors,
		RawScore:      rawScore,
		Amplification: amplification,
		FinalScore:    finalScore,
		Level:         scoreToLevel(finalScore, config),
		Recommended:   scoreToDecision(finalScore, config),
	}
}

func Tier(action string, policyRiskOverride *kerneldomain.RiskLevel, config Config) kerneldomain.RiskLevel {
	if policyRiskOverride != nil {
		switch *policyRiskOverride {
		case kerneldomain.RiskHigh:
			return kerneldomain.RiskHigh
		case kerneldomain.RiskMedium:
			return kerneldomain.RiskMedium
		case kerneldomain.RiskLow:
			return kerneldomain.RiskLow
		}
	}
	action = strings.TrimSpace(action)
	if slices.Contains(config.HighActions, action) {
		return kerneldomain.RiskHigh
	}
	if slices.Contains(config.MediumActions, action) {
		return kerneldomain.RiskMedium
	}
	return kerneldomain.RiskLow
}

func DecideFromPolicy(effect kerneldomain.Decision, tier kerneldomain.RiskLevel) (kerneldomain.Decision, bool) {
	switch effect {
	case kerneldomain.DecisionDeny:
		return kerneldomain.DecisionDeny, true
	case kerneldomain.DecisionRequireApproval:
		return kerneldomain.DecisionRequireApproval, true
	case kerneldomain.DecisionAllow:
		if tier == kerneldomain.RiskHigh {
			return kerneldomain.DecisionRequireApproval, true
		}
		return kerneldomain.DecisionAllow, true
	default:
		return "", false
	}
}

func DefaultDecision(tier kerneldomain.RiskLevel) kerneldomain.Decision {
	if tier == kerneldomain.RiskHigh {
		return kerneldomain.DecisionRequireApproval
	}
	return kerneldomain.DecisionAllow
}

func evaluateFactors(request kerneldomain.Request, history History, config Config, now time.Time) []kerneldomain.RiskFactor {
	factors := make([]kerneldomain.RiskFactor, 0, 6)
	factors = append(factors, actionFactor(request.Action, config))
	factors = append(factors, offHoursFactor(now, config))
	factors = append(factors, actorHistoryFactor(history.ActorHistory, config))
	factors = append(factors, frequencyFactor(history.RecentFrequency, config))
	factors = append(factors, successRateFactor(history.SuccessRate, config))
	factors = append(factors, targetFactor(request.Target.System, config))
	return factors
}

func actionFactor(action string, config Config) kerneldomain.RiskFactor {
	action = strings.TrimSpace(action)
	switch {
	case slices.Contains(config.HighActions, action):
		return kerneldomain.RiskFactor{Name: "action_type", Score: 0.4, Active: true, Reason: action + " is high-risk action"}
	case slices.Contains(config.MediumActions, action):
		return kerneldomain.RiskFactor{Name: "action_type", Score: 0.2, Active: true, Reason: action + " is medium-risk action"}
	default:
		return kerneldomain.RiskFactor{Name: "action_type", Score: 0.1, Reason: action + " is low-risk action"}
	}
}

func offHoursFactor(now time.Time, config Config) kerneldomain.RiskFactor {
	hour := now.UTC().Hour()
	if hour < config.BusinessHours.Start || hour >= config.BusinessHours.End {
		return kerneldomain.RiskFactor{Name: "off_hours", Score: 0.2, Active: true, Reason: "request at off-hours"}
	}
	return kerneldomain.RiskFactor{Name: "off_hours"}
}

func actorHistoryFactor(history int, config Config) kerneldomain.RiskFactor {
	switch {
	case history <= config.ActorThresholds.Unknown:
		return kerneldomain.RiskFactor{Name: "actor_unknown", Score: 0.3, Active: true, Reason: "unknown actor, no previous requests"}
	case history < config.ActorThresholds.New:
		return kerneldomain.RiskFactor{Name: "actor_unknown", Score: 0.15, Active: true, Reason: "new actor with limited history"}
	default:
		return kerneldomain.RiskFactor{Name: "actor_unknown"}
	}
}

func frequencyFactor(count int, config Config) kerneldomain.RiskFactor {
	switch {
	case count > config.FrequencyThresholds.Critical:
		return kerneldomain.RiskFactor{Name: "frequency_anomaly", Score: 0.3, Active: true, Reason: "frequency above critical threshold"}
	case count > config.FrequencyThresholds.Warning:
		return kerneldomain.RiskFactor{Name: "frequency_anomaly", Score: 0.15, Active: true, Reason: "frequency above warning threshold"}
	default:
		return kerneldomain.RiskFactor{Name: "frequency_anomaly"}
	}
}

func successRateFactor(value float64, config Config) kerneldomain.RiskFactor {
	if value < 0 {
		return kerneldomain.RiskFactor{Name: "execution_history"}
	}
	switch {
	case value < config.SuccessRateThresholds.Low:
		return kerneldomain.RiskFactor{Name: "execution_history", Score: 0.3, Active: true, Reason: "low historical success rate"}
	case value < config.SuccessRateThresholds.Moderate:
		return kerneldomain.RiskFactor{Name: "execution_history", Score: 0.1, Active: true, Reason: "moderate historical success rate"}
	case value >= config.SuccessRateThresholds.Excellent:
		return kerneldomain.RiskFactor{Name: "execution_history", Score: -0.15, Reason: "excellent historical success rate"}
	default:
		return kerneldomain.RiskFactor{Name: "execution_history"}
	}
}

func targetFactor(system string, config Config) kerneldomain.RiskFactor {
	system = strings.TrimSpace(strings.ToLower(system))
	if slices.Contains(config.SensitiveSystems, system) {
		return kerneldomain.RiskFactor{Name: "target_sensitivity", Score: 0.3, Active: true, Reason: "target is sensitive system"}
	}
	return kerneldomain.RiskFactor{Name: "target_sensitivity"}
}

func calculateAmplification(factors []kerneldomain.RiskFactor, config Config) float64 {
	active := make(map[string]bool, len(factors))
	count := 0
	for _, factor := range factors {
		if factor.Active {
			active[factor.Name] = true
			count++
		}
	}

	amplification := 1.0
	for _, rule := range config.Amplifications {
		if allFactorsActive(rule.Factors, active) {
			amplification = max(amplification, rule.Multiplier)
		}
	}
	if count >= 4 {
		amplification = max(amplification, 2.5)
	}
	return min(amplification, config.Thresholds.MaxAmplification)
}

func allFactorsActive(factors []string, active map[string]bool) bool {
	for _, name := range factors {
		if !active[name] {
			return false
		}
	}
	return true
}

func applyPolicyOverride(override kerneldomain.RiskLevel, currentScore float64, config Config) float64 {
	switch override {
	case kerneldomain.RiskHigh:
		if currentScore < config.Thresholds.RequireApproval {
			return config.Thresholds.RequireApproval
		}
	case kerneldomain.RiskMedium:
		if currentScore < config.Thresholds.EnhancedLog {
			return config.Thresholds.EnhancedLog
		}
	case kerneldomain.RiskLow:
		if currentScore > config.Thresholds.Allow {
			return config.Thresholds.Allow * 0.9
		}
	}
	return currentScore
}

func scoreToLevel(score float64, config Config) kerneldomain.RiskLevel {
	switch {
	case score >= config.Thresholds.RequireApproval:
		return kerneldomain.RiskHigh
	case score >= config.Thresholds.EnhancedLog:
		return kerneldomain.RiskMedium
	default:
		return kerneldomain.RiskLow
	}
}

func scoreToDecision(score float64, config Config) kerneldomain.Decision {
	switch {
	case score >= config.Thresholds.Deny:
		return kerneldomain.DecisionDeny
	case score >= config.Thresholds.RequireApproval:
		return kerneldomain.DecisionRequireApproval
	default:
		return kerneldomain.DecisionAllow
	}
}

func sumFactors(factors []kerneldomain.RiskFactor) float64 {
	total := 0.0
	for _, factor := range factors {
		total += factor.Score
	}
	if total < 0 {
		return 0
	}
	return total
}

func normalizeConfig(config Config) Config {
	defaults := DefaultConfig()
	if len(config.HighActions) == 0 {
		config.HighActions = defaults.HighActions
	}
	if len(config.MediumActions) == 0 {
		config.MediumActions = defaults.MediumActions
	}
	if config.BusinessHours == (BusinessHours{}) {
		config.BusinessHours = defaults.BusinessHours
	}
	if config.FrequencyThresholds == (FrequencyThresholds{}) {
		config.FrequencyThresholds = defaults.FrequencyThresholds
	}
	if config.ActorThresholds == (ActorThresholds{}) {
		config.ActorThresholds = defaults.ActorThresholds
	}
	if config.SuccessRateThresholds == (SuccessRateThresholds{}) {
		config.SuccessRateThresholds = defaults.SuccessRateThresholds
	}
	if len(config.Amplifications) == 0 {
		config.Amplifications = defaults.Amplifications
	}
	if len(config.SensitiveSystems) == 0 {
		config.SensitiveSystems = defaults.SensitiveSystems
	}
	if config.Thresholds == (Thresholds{}) {
		config.Thresholds = defaults.Thresholds
	}
	if config.Thresholds.Allow == 0 {
		config.Thresholds.Allow = thresholdAllow
	}
	if config.Thresholds.EnhancedLog == 0 {
		config.Thresholds.EnhancedLog = thresholdEnhancedLog
	}
	if config.Thresholds.RequireApproval == 0 {
		config.Thresholds.RequireApproval = defaults.Thresholds.RequireApproval
	}
	if config.Thresholds.Deny == 0 {
		config.Thresholds.Deny = defaults.Thresholds.Deny
	}
	if config.Thresholds.MaxAmplification == 0 {
		config.Thresholds.MaxAmplification = defaults.Thresholds.MaxAmplification
	}
	return config
}
