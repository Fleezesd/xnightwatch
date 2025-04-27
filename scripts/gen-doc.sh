#!/usr/bin/env bash

for top in pkg
do
  for d in $(find $top -type d)
  do
    if [[ "$d" =~ "pkg/generated" ]];then
      continue
    fi

    if [ ! -f $d/doc.go ]; then
      if ls $d/*.go > /dev/null 2>&1; then
        echo $d/doc.go
        echo "package $(basename $d) // import \"github.com/fleezesd/xnightwatch/$d\"" > $d/doc.go
      fi
    fi
  done
done