package prom

import (
	"fmt"
	"testing"

	"github.com/prometheus/prometheus/model/labels"
)

func TestInject(t *testing.T) {
	fmt.Printf(InjectLabel("sum(cpu_usage_active{instance=\"111\"})", "asset_id", "20", labels.MatchEqual))
}
