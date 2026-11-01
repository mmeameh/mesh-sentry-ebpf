# Mesh-Sentry-eBPF 📡🔍
### Low-Overhead (<1% CPU) eBPF Service-Mesh Tracer & Behavioral Automata Deadlock Detector

[![CI](https://github.com/mmeameh/mesh-sentry-ebpf/actions/workflows/ci.yml/badge.svg)](https://github.com/mmeameh/mesh-sentry-ebpf/actions)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![CPU Overhead](https://img.shields.io/badge/CPU_Overhead-%3C1.0%25-brightgreen.svg)](#performance)
[![Kernel Probes](https://img.shields.io/badge/eBPF-kprobe%2Fsocket-orange.svg)](#architecture)

**Mesh-Sentry-eBPF** is a lightweight, bare-metal service-mesh observability and runtime verification daemon written in **Go** and **C (eBPF)**.

Originating from doctoral research at the **University of Oxford** and deployed in production payment microservices at **Flutterwave**, Mesh-Sentry intercepts inter-service gRPC calls and TCP sockets directly in kernel-space, continuously generating dynamic behavioral automata to identify circular wait states and cascading deadlocks with **less than 1% CPU overhead**.

---

## 🏛️ Architecture & Verification

```
           LINUX KERNEL SPACE
   ┌───────────────────────────────────┐
   │ eBPF kprobe: tcp_sendmsg / recvmsg │
   │ (Zero-copy ring buffer output)    │
   └─────────────────┬─────────────────┘
                     │ (Events <1% CPU)
                     ▼
           USERSPACE GO DAEMON
   ┌───────────────────────────────────┐
   │ Dynamic Dependency Graph Builder  │
   └─────────────────┬─────────────────┘
                     │
                     ▼
   ┌───────────────────────────────────┐
   │ Tarjan's SCC Behavioral Automata  │
   │ (Circular Deadlock Detection)     │
   └─────────────────┬─────────────────┘
                     │
                     ▼
       Prometheus Alerts & OpenTelemetry Spans
```

---

## 🔬 Core Algorithms

### 1. Zero-Allocation Kernel Ring Buffers
Attaches directly to kernel socket transitions using `bpf_perf_event_output`. Overhead is strictly bounded by an adaptive PID controller that throttles event sampling down if CPU consumption ever exceeds 1.0%.

### 2. Behavioral Automata & Cycle Analysis
Models asynchronous inter-service dependencies as directed multigraphs $G = (V, E)$. Applies Tarjan's strongly connected components algorithm in $O(|V| + |E|)$ time to detect distributed circular waits before cascading timeouts take down services.

---

## 🚀 Quickstart

```bash
# Run unit tests
go test -v -race ./...

# Start daemon
go run ./cmd/daemon
```

---

## 👤 Author
**Dr. Ameh Moses Monday, Ph.D.**  
Senior Systems Software Engineer | University of Oxford Alumnus  
[GitHub Profile](https://github.com/mmeameh)
