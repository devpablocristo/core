package policy

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/devpablocristo/core/governance/go/domain"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

type Evaluator struct {
	env    *cel.Env
	envErr error
	mu     sync.Mutex
	progs  map[string]cel.Program
}

func NewEvaluator() *Evaluator {
	env, err := cel.NewEnv(
		cel.Variable("request", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("time", cel.MapType(cel.StringType, cel.DynType)),
	)
	if err != nil {
		return &Evaluator{envErr: err, progs: make(map[string]cel.Program)}
	}
	return &Evaluator{env: env, progs: make(map[string]cel.Program)}
}

func (e *Evaluator) Match(request domain.Request, item domain.Policy, now time.Time) (bool, error) {
	if strings.TrimSpace(item.ActionFilter) != "" && strings.TrimSpace(item.ActionFilter) != strings.TrimSpace(request.Action) {
		return false, nil
	}
	if strings.TrimSpace(item.SystemFilter) != "" && strings.TrimSpace(item.SystemFilter) != strings.TrimSpace(request.Target.System) {
		return false, nil
	}
	return e.Matches(item.Expression, request, now)
}

func (e *Evaluator) Matches(expression string, request domain.Request, now time.Time) (bool, error) {
	if strings.TrimSpace(expression) == "" {
		return true, nil
	}
	prog, err := e.program(expression)
	if err != nil {
		return false, err
	}

	result, _, err := prog.Eval(map[string]any{
		"request": RequestToMap(request),
		"time": map[string]any{
			"hour":        now.UTC().Hour(),
			"day_of_week": int(now.UTC().Weekday()),
		},
	})
	if err != nil {
		return false, fmt.Errorf("eval policy: %w", err)
	}
	if result.Type() != types.BoolType {
		return false, fmt.Errorf("policy must return bool, got %s", result.Type())
	}
	value, ok := result.Value().(bool)
	if !ok {
		return false, fmt.Errorf("policy result is not bool")
	}
	return value, nil
}

func RequestToMap(request domain.Request) map[string]any {
	return map[string]any{
		"id": request.ID,
		"subject": map[string]any{
			"type": string(request.Subject.Type),
			"id":   request.Subject.ID,
			"name": request.Subject.Name,
		},
		"action": request.Action,
		"target": map[string]any{
			"system":   request.Target.System,
			"resource": request.Target.Resource,
		},
		"params":     cloneMap(request.Params),
		"metadata":   cloneMap(request.Metadata),
		"reason":     request.Reason,
		"context":    request.Context,
		"created_at": request.CreatedAt.UTC().Format(time.RFC3339),
	}
}

func cloneMap(values map[string]any) map[string]any {
	if len(values) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(values))
	for key, value := range values {
		out[key] = value
	}
	return out
}

func (e *Evaluator) program(expression string) (cel.Program, error) {
	if e.envErr != nil {
		return nil, e.envErr
	}
	e.mu.Lock()
	if prog, ok := e.progs[expression]; ok {
		e.mu.Unlock()
		return prog, nil
	}
	e.mu.Unlock()

	ast, issues := e.env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, issues.Err()
	}
	if ast.OutputType() != cel.BoolType {
		return nil, fmt.Errorf("expression must return bool")
	}
	prog, err := e.env.Program(ast)
	if err != nil {
		return nil, err
	}

	e.mu.Lock()
	e.progs[expression] = prog
	e.mu.Unlock()
	return prog, nil
}
