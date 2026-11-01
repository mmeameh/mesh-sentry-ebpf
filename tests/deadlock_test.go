package tests

import (
	"testing"

	"github.com/mmeameh/mesh-sentry-ebpf/pkg/automata"
	"github.com/mmeameh/mesh-sentry-ebpf/pkg/throttle"
)

func TestDeadlockDetectionCircularWait(t *testing.T) {
	graph := automata.NewDependencyGraph()

	// Construct acyclic call tree: OrderService -> InventoryService -> Database
	graph.AddDependency("OrderService", "InventoryService")
	graph.AddDependency("InventoryService", "Database")

	detector := automata.NewDeadlockDetector()
	cycles := detector.FindCycles(graph)
	if len(cycles) != 0 {
		t.Fatalf("expected 0 deadlocks in acyclic tree, found %d: %v", len(cycles), cycles)
	}

	// Introduce circular dependency: Database -> PaymentService -> OrderService
	graph.AddDependency("Database", "PaymentService")
	graph.AddDependency("PaymentService", "OrderService")

	cycles = detector.FindCycles(graph)
	if len(cycles) != 1 {
		t.Fatalf("expected 1 deadlock cycle detected, found %d", len(cycles))
	}

	cycleNodes := cycles[0]
	if len(cycleNodes) != 4 {
		t.Errorf("expected 4 services involved in circular wait, got %d", len(cycleNodes))
	}
}

func TestOverheadThrottlingUnderOnePercent(t *testing.T) {
	ctrl := throttle.NewOverheadController(1.0) // 1.0% max CPU

	// Baseline normal usage (0.3% CPU)
	ctrl.UpdateUsage(0.3)
	if ctrl.SampleRate() < 0.9 {
		t.Errorf("expected near full sampling at 0.3%% CPU, got %f", ctrl.SampleRate())
	}

	// CPU spike (1.8% CPU > 1.0% budget)
	ctrl.UpdateUsage(1.8)
	if ctrl.SampleRate() >= 1.0 {
		t.Errorf("expected throttle down when CPU exceeds budget, got %f", ctrl.SampleRate())
	}
}
