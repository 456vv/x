@echo off

set GOEXPERIMENT=noregabi
go mod tidy

if "%1" == "" goto input

for /f "delims=" %%I in (%1) do (go get %%I && qexp -outdir . -addtags "//go:build igop_lib;// +build igop_lib" %%I)
goto exit

:input
	set /p flag=enter pkg path:
	if "%flag%" == "" goto all
	go get -u %flag%
	qexp -outdir . -addtags "//go:build igop_lib;// +build igop_lib" %flag%
	pause
	goto input
	exit 0
:all
	for /f "tokens=1* delims=:" %%a in ('findstr /n .* igop_pkg_list.txt') do (
		if "%%b" == "" (pause) else (
			if not exist %%b (echo %%b && go get %%b && qexp -outdir . -addtags "//go:build igop_lib;// +build igop_lib" %%b) else (echo skip %%b)
		)
	)
:exit0
echo �������
pause
exit 0

