#!/usr/bin/env bash
# Runs your tests against deliberately buggy copies of Quote (the mutant_*.go
# files). A good test suite fails on every one of them: each mutant must be
# "caught". You don't need to understand build tags to use this.
set -u
cd "$(dirname "$0")"

if ! go test -count=1 . >/dev/null 2>&1; then
	echo "Your tests fail against the real Quote. Fix that first: go test ./problems/11-table-tests/..."
	exit 1
fi

caught=0
total=0
for f in mutant_*.go; do
	tag="${f%.go}"
	total=$((total + 1))
	if go test -count=1 -tags "mutant,$tag" . >/dev/null 2>&1; then
		echo "SURVIVED  $tag   <- every test passed against this bug; add a case that catches it"
	else
		echo "caught    $tag"
		caught=$((caught + 1))
	fi
done

echo "$caught/$total mutants caught"
[ "$caught" -eq "$total" ]
