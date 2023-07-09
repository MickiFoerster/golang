// +build ignore
#include "hello.bpf.h"

#include <linux/types.h>

struct bpf_map_def SEC("maps") mymap = {.type = BPF_MAP_TYPE_PERF_EVENT_ARRAY,
                                        .key_size = sizeof(int),
                                        .value_size = sizeof(__u32),
                                        .max_entries = 1024};

SEC("kprobe/sys_execve")

int hello(void *ctx) {
    __u64 data = 42;
    bpf_perf_event_output(ctx, &mymap, BPF_F_CURRENT_CPU, &data, sizeof(data));

    return 0;
}
