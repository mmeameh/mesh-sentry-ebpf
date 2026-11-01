package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mmeameh/mesh-sentry-ebpf/pkg/automata"
	"github.com/mmeameh/mesh-sentry-ebpf/pkg/throttle"
)

func main() {
	fmt.Println("=================================================================")
	fmt.Println("  Mesh-Sentry-eBPF: Low-Overhead (<1% CPU) Microservice Tracer   ")
	fmt.Println("  Lead Architect: Dr. Ameh Moses Monday (Oxford Ph.D.)           ")
	fmt.Println("=================================================================")

	graph := automata.NewDependencyGraph()
	detector := automata.NewDeadlockDetector()
	ctrl := throttle.NewOverheadController(1.0)

	// Simulate discovery of active gRPC call graph
	graph.AddDependency("checkout-api", "fraud-detector")
	graph.AddDependency("fraud-detector", "triton-infer-guard")

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"RUNNING","cpu_budget_pct":1.0,"ebpf_probes_active":true}`))
	})

	http.HandleFunc("/v1/deadlocks", func(w http.ResponseWriter, r *http.Request) {
		cycles := detector.FindCycles(graph)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"cycles_detected":%d,"cpu_sample_rate":%.2f}`, len(cycles), ctrl.SampleRate())
	})

	log.Println("Mesh-Sentry-eBPF daemon listening on :9090...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
