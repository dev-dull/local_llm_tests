#!/usr/bin/env bash
#
# Scaffold the results directory structure for a model under test.
#
# Prompts for a model name and a number of runs, then for every test in
# prompts/ creates results/<model>/<test-name>/<n>/ and copies the test's
# prompt files into each numbered directory. If numbered runs already exist,
# new ones continue after the highest existing number.

set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
prompts_dir="$root/prompts"
results_dir="$root/results"

if [ ! -d "$prompts_dir" ]; then
    echo "error: no prompts directory at $prompts_dir" >&2
    exit 1
fi

read -r -p "Model name: " model
if [ -z "$model" ] || [[ "$model" == */* ]]; then
    echo "error: model name must be non-empty and contain no slashes" >&2
    exit 1
fi

read -r -p "Number of runs: " runs
if ! [[ "$runs" =~ ^[1-9][0-9]*$ ]]; then
    echo "error: number of runs must be a positive integer" >&2
    exit 1
fi

found_tests=0
for prompt_dir in "$prompts_dir"/*/; do
    [ -d "$prompt_dir" ] || continue
    found_tests=1
    test_name="$(basename "$prompt_dir")"
    test_dir="$results_dir/$model/$test_name"

    # Continue numbering after the highest existing run, if any.
    start=1
    if [ -d "$test_dir" ]; then
        for existing in "$test_dir"/*/; do
            [ -d "$existing" ] || continue
            n="$(basename "$existing")"
            [[ "$n" =~ ^[0-9]+$ ]] || continue
            if [ "$n" -ge "$start" ]; then
                start=$((n + 1))
            fi
        done
    fi

    for ((i = start; i < start + runs; i++)); do
        run_dir="$test_dir/$i"
        mkdir -p "$run_dir"
        cp -R "$prompt_dir"/. "$run_dir/"
        echo "created $run_dir"
    done
done

if [ "$found_tests" -eq 0 ]; then
    echo "error: no test directories found in $prompts_dir" >&2
    exit 1
fi
