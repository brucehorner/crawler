#!/bin/bash

coveragefile="coverage.txt"
[ -e $coveragefile ] && rm $coveragefile
go clean

DEPS=$(go list -f '{{join .Deps "\n"}}' | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}')
printf "Dependencies are:\n"
echo $DEPS | sed "s/ \|^/\n\t/g"

go get -v $DEPS
go fmt -x *.go
go build
result=$?

if [ $result -eq 0 ]
then
  go test -coverprofile=$coveragefile
  result=$?
fi

if [ $result -eq 0 ]
then
  go install
  result=$?
fi

exit $result

