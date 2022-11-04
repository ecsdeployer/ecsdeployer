#!/bin/sh
set -e

SED="sed"
if which gsed >/dev/null 2>&1; then
	SED="gsed"
fi

go generate ./www/
go run . schema -o ./www/docs/static/schema.json

rm -rf www/docs/cmd/*.md
go run . docs
"$SED" \
	-i'' \
	-e 's/SEE ALSO/See also/g' \
	-e 's/Options inherited from parent commands/Global Options/g' \
	-e 's/^## /# /g' \
	-e 's/^### /## /g' \
	-e 's/^#### /### /g' \
	-e 's/^##### /#### /g' \
	-e 's/^# ecsdeployer/---\nhide:\n  toc: true\n---\n# ecsdeployer/g' \
	-e 's/^\* \[ecsdeployer/* [`ecsdeployer/g' \
	-e 's/](ecsdeployer/`](ecsdeployer/g' \
	./www/docs/cmd/*.md
