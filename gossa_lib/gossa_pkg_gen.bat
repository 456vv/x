@echo off

set GOEXPERIMENT=noregabi
go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
go mod tidy -compat=1.17

if "%1" == "" goto input

for /f "delims=" %%I in (%1) do (gossa -outdir . -addtags "//go:build gossa_lib;// +build gossa_lib" %%I)
goto exit

:input
	set /p flag=enter pkg path:
	if "%flag%" == "" goto all
	gossa -outdir . -addtags "//go:build gossa_lib;// +build gossa_lib" %flag%
	pause
	goto input
	exit 0
:all
	for /f "tokens=1* delims=:" %%a in ('findstr /n .* gossa_pkg_list.txt') do (
		if "%%b" == "" (pause) else (
			if not exist %%b (echo %%b && gossa -outdir . -addtags "//go:build gossa_lib;// +build gossa_lib" %%b) else (echo skip %%b)
		)
	)
:exit0
echo 生成完成
pause
exit 0

