#!/usr/bin/env ruby

lines = File.readlines 'CHANGELOG.md'

# Collate lines by numbered version
versions = lines.slice_when { |_, line| line =~ /^## \[\d/ }.to_a

# Get lines from the most recent version
current = versions[1].map { |line| line.chomp '' }

# Replace the header
current[0] = "## Changelog\n"

puts current.join("\n")
