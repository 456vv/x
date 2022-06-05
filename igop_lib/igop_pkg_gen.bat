@echo off

set GOEXPERIMENT=noregabi
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy -compat=1.17

if "%1" == "" goto input

for /f "delims=" %%I in (%1) do (qexp -outdir . -addtags "//go:build igop_lib;// +build igop_lib" %%I)
goto exit

:input
	set /p flag=enter pkg path:
	if "%flag%" == "" goto all
	qexp -outdir . -addtags "//go:build igop_lib;// +build igop_lib" %flag%
	pause
	goto input
	exit 0
:all
	for /f "tokens=1* delims=:" %%a in ('findstr /n .* igop_pkg_list.txt') do (
		if "%%b" == "" (pause) else (
			if not exist %%b (echo %%b && qexp -outdir . -addtags "//go:build igop_lib;// +build igop_lib" %%b) else (echo skip %%b)
		)
	)
:exit0
echo 生成完成
pause
exit 0

