module github.com/456vv/x/goplus_lib/cmd/main

go 1.16

replace github.com/456vv/viot/v2 => ../../../../viot

replace github.com/golang/freetype => ../../../../../../github.com/golang/freetype

require (
	github.com/456vv/vbody v1.2.3
	github.com/456vv/vcipher v1.0.0
	github.com/456vv/vconn v1.1.2
	github.com/456vv/vconnpool/v2 v2.1.4
	github.com/456vv/verifycode v1.0.3
	github.com/456vv/verror v1.1.0 // indirect
	github.com/456vv/vforward v1.0.8
	github.com/456vv/viot/v2 v2.0.0-00010101000000-000000000000
	github.com/456vv/vmap/v2 v2.3.1
	github.com/456vv/vweb/v2 v2.4.6
	github.com/456vv/x/db v0.0.0-20210530134931-1dc04f4064f0
	github.com/456vv/x/smtp v0.0.0-20210530134931-1dc04f4064f0
	github.com/456vv/x/sqltable v0.0.0-20210530134931-1dc04f4064f0
	github.com/456vv/x/ticker v0.0.0-20210530134931-1dc04f4064f0
	github.com/456vv/x/watch v0.0.0-20210530134931-1dc04f4064f0
	github.com/88250/lute v1.7.2
	github.com/bep/golibsass v1.0.0
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/freetype v0.0.0-00010101000000-000000000000 // indirect
	github.com/goplus/reflectx v0.5.2
	github.com/lib/pq v1.10.2 // indirect
	github.com/mattn/go-sqlite3 v1.14.7 // indirect
	github.com/pelletier/go-toml v1.9.1
	golang.org/x/image v0.0.0-20210504121937-7319ad40d33e // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
