package prom

import (
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func injectLabels(expr parser.Expr, match labels.MatchType, name, value string) {
	switch e := expr.(type) {
	case *parser.AggregateExpr:
		injectLabels(e.Expr, match, name, value)
	case *parser.Call:
		for _, v := range e.Args {
			injectLabels(v, match, name, value)
		}
	case *parser.ParenExpr:
		injectLabels(e.Expr, match, name, value)
	case *parser.UnaryExpr:
		injectLabels(e.Expr, match, name, value)
	case *parser.BinaryExpr:
		injectLabels(e.LHS, match, name, value)
		injectLabels(e.RHS, match, name, value)
	case *parser.VectorSelector:
		l := genMetricLabel(match, name, value)
		e.LabelMatchers = append(e.LabelMatchers, l)
		return
	case *parser.MatrixSelector:
		injectLabels(e.VectorSelector, match, name, value)
	case *parser.SubqueryExpr:
		injectLabels(e.Expr, match, name, value)
	case *parser.NumberLiteral, *parser.StringLiteral:
		return
	default:
		panic(errors.Errorf("unhandled expression of type: %T", expr))
	}
	return
}

func genMetricLabel(match labels.MatchType, name, value string) *labels.Matcher {
	m, err := labels.NewMatcher(match, name, value)
	if nil != err {
		return nil
	}

	return m
}

func InjectLabel(promql string, label string, value string, op labels.MatchType) string {
	expr, _ := parser.ParseExpr(promql)
	injectLabels(expr, labels.MatchEqual, label, value)
	return expr.String()
}
