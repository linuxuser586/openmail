#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

year=$(date +'%Y')

files_without_header=()

header_match_lines=("^// Copyright $year The OpenMail Authors$" "^// SPDX-License-Identifier: Apache-2.0$" "")

# Retrieve the list of newly added files
mapfile -t newly_added_files < <(git diff --name-only --diff-filter=A --cached -- '*.go')
if [ "${#newly_added_files[@]}" != 0 ]
then
    # Check for Copyright statement
    for newly_added_file in "${newly_added_files[@]}"
    do
        for i in "${!header_match_lines[@]}"
        do
            mapfile -t txt < <(head "-n$((i+1))" "$newly_added_file" | tail -n1 | grep -q "${header_match_lines[$i]}" || echo "not found")
            if [ "${#txt[@]}" != 0 ]
            then
                files_without_header+=("$newly_added_file")
                break
            fi
        done
    done

    if [ "${#files_without_header[@]}" != 0  ]
    then
        # License header missing in files
        for line in "${header_match_lines[@]}"
        do
            line="${line#"^"}"
            echo "${line%"\$"}"
        done
        echo "The above license header was not found in the following new files:"
        for file in "${files_without_header[@]}"
        do
            :
            echo "   - $file";
        done
        exit 1;
    else
        echo "All new go files have the correct license header.";
        exit 0;
    fi
else
    echo "No new go files"
fi