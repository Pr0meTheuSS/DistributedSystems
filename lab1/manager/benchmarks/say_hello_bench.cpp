#include <userver/utest/using_namespace_userver.hpp>

#include <string_view>

#include <benchmark/benchmark.h>
#include <userver/engine/run_standalone.hpp>

void HelloBenchmark(benchmark::State& state)
{
    engine::RunStandalone([&] {
        for (auto _ : state) {
        }
    });
}

BENCHMARK(HelloBenchmark);