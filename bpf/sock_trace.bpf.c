// +build ignore

#include <linux/bpf.h>
#include <linux/ptrace.h>
#include <linux/tcp.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>

struct socket_event_t {
    __u32 pid;
    __u32 saddr;
    __u32 daddr;
    __u16 sport;
    __u16 dport;
    __u64 latency_ns;
    __u64 bytes_sent;
};

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(__u32));
    __uint(value_size, sizeof(__u32));
} events SEC(".maps");

SEC("kprobe/tcp_sendmsg")
int BPF_KPROBE(trace_tcp_sendmsg, struct sock *sk, struct msghdr *msg, size_t size) {
    struct socket_event_t event = {};
    event.pid = bpf_get_current_pid_tgid() >> 32;
    event.bytes_sent = size;
    event.latency_ns = bpf_ktime_get_ns();

    // Export non-blocking event to userspace ring buffer (<1% CPU)
    bpf_perf_event_output(ctx, &events, BPF_F_CURRENT_CPU, &event, sizeof(event));
    return 0;
}

char LICENSE[] SEC("license") = "Dual BSD/GPL";
