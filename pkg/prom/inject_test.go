package prom

import (
	"fmt"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
)

func TestInject(t *testing.T) {
	fmt.Println(InjectLabel("sum(cpu_usage_active)", "asset_id", "20", labels.MatchNotEqual))
}
