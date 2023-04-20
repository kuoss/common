#!/bin/bash
MIN_COVER=80

cd $(dirname $0)/..
export PS4='[$(basename $0):$LINENO] '

echo + go test ./... -v -failfast -race -covermode=atomic -coverprofile /tmp/cover.out
       go test ./... -v -failfast -race -covermode=atomic -coverprofile /tmp/cover.out
if [[ $? != 0 ]]; then
    echo "❌ FAIL - test failed"
    exit 1
fi

COVER=$(go tool cover -func /tmp/cover.out | tail -1 | grep -oP [0-9]+ | head -1)
rm -f /tmp/cover.out
if [ $COVER -lt $MIN_COVER ]; then
    echo "⚠️ WARN - total coverage: ${COVER}% (<${MIN_COVER}%)"
    exit
fi
echo "✔️ OK - total coverage: ${COVER}% (>=${MIN_COVER}%)"
