#!/usr/bin/env bash
#
# Scaffold the results directory structure for a model under test.
#
# Prompts for a model name, then for a number of runs per test in prompts/
# (0 or blank skips a test). For each selected test it creates
# results/<model>/<test-name>/<n>/ and copies the test's prompt files into
# each numbered directory.
#
# This script is strictly non-destructive: it never deletes anything and
# never overwrites an existing file. Every directory and file is checked
# before it is created or copied; anything that already exists is kept as-is
# (with a "kept" notice) and scaffolding continues around it. Numbering
# always continues after the highest existing run.

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

# Guard against typos silently creating a fresh model tree.
if [ ! -d "$results_dir/$model" ]; then
    echo "note: results/$model does not exist yet. Existing models:"
    for d in "$results_dir"/*/; do
        [ -d "$d" ] && echo "  - $(basename "$d")"
    done
    read -r -p "Create new model directory results/$model? [y/N] " answer
    case "$answer" in
        y|Y|yes|YES) ;;
        *) echo "aborted"; exit 1 ;;
    esac
fi

# Copy every file from $1 (a prompt dir) into $2 (a run dir), creating
# directories as needed and never overwriting an existing file.
copy_prompt_files() {
    local src="$1" dest_dir="$2" f rel dest
    while IFS= read -r -d '' f; do
        rel="${f#"$src"/}"
        dest="$dest_dir/$rel"
        if [ -e "$dest" ]; then
            echo "kept existing $dest (not overwritten)"
            continue
        fi
        mkdir -p "$(dirname "$dest")"
        cp "$f" "$dest"
    done < <(find "$src" -type f -print0)
}

found_tests=0
for prompt_dir in "$prompts_dir"/*/; do
    [ -d "$prompt_dir" ] || continue
    prompt_dir="${prompt_dir%/}"
    found_tests=1
    test_name="$(basename "$prompt_dir")"
    test_dir="$results_dir/$model/$test_name"

    # Continue numbering after the highest existing run, if any.
    start=1
    existing_runs=0
    if [ -d "$test_dir" ]; then
        for existing in "$test_dir"/*/; do
            [ -d "$existing" ] || continue
            n="$(basename "$existing")"
            [[ "$n" =~ ^[0-9]+$ ]] || continue
            existing_runs=$((existing_runs + 1))
            if [ "$n" -ge "$start" ]; then
                start=$((n + 1))
            fi
        done
    fi

    read -r -p "Runs to add for $test_name ($existing_runs existing) [0 to skip]: " runs
    if [ -z "$runs" ] || [ "$runs" = "0" ]; then
        echo "skipped $test_name"
        continue
    fi
    if ! [[ "$runs" =~ ^[1-9][0-9]*$ ]]; then
        echo "error: number of runs must be a positive integer (or 0 to skip)" >&2
        exit 1
    fi

    for ((i = start; i < start + runs; i++)); do
        run_dir="$test_dir/$i"
        if [ -e "$run_dir" ]; then
            # Shouldn't happen (numbering starts past the highest existing
            # run), but if it does: add only missing files, overwrite nothing.
            echo "exists $run_dir — filling in missing files only"
        else
            mkdir -p "$run_dir"
            echo "created $run_dir"
        fi
        copy_prompt_files "$prompt_dir" "$run_dir"
    done
done

if [ "$found_tests" -eq 0 ]; then
    echo "error: no test directories found in $prompts_dir" >&2
    exit 1
fi
