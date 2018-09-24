# gccgo-export-dumper

Reads objects/archives built by gccgo, dumps export data.

Example:

```
% cd $GOPATH/foo
% go build -compiler gccgo .
% gccgo-export-dumper $GOPATH/pkg/gccgo_linux_amd64/libfoo.a
v2;
package foo;
pkgpath foo;
import fmt fmt "fmt";
import runtime runtime "runtime";
init fmt fmt..import poll internal_zpoll..import testlog internal_ztestlog..import io io..import os os..import reflect reflect..import runtime runtime..import sys runtime_zinternal_zsys..import strconv strconv..import sync sync..import syscall syscall..import time time..import unicode unicode..import;
init_graph 0 1 0 2 0 3 0 4 0 5 0 6 0 7 0 8 0 9 0 10 0 11 0 12 1 3 1 6 1 7 1 9 1 10 1 11 3 6 3 7 3 9 4 1 4 2 4 3 4 6 4 7 4 9 4 10 4 11 5 6 5 7 5 8 5 9 5 12 6 7 8 6 8 7 9 6 9 7 10 6 10 7 10 9 11 6 11 7 11 9 11 10;
func Something ();
checksum 2B8A4D3A17552141115613BAC882E00DF8B352A4;
```



