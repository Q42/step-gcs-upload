#!/bin/bash

THIS_SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Start go program
cd "${THIS_SCRIPTDIR}"

go get -u golang.org/x/net/context
go get -u golang.org/x/oauth2/google
go get -u google.golang.org/api/storage/...
go run ./program.go
ex_code=$?

if [ ${ex_code} -eq 0 ] ; then
  exit 0
fi

exit 1