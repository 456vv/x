@echo off

go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct
go mod tidy

if "%1" == "" goto input

for /f "delims=" %%I in (%1) do (yaegi extract -tag yaegi_lib %%I)
goto exit

:input
	set /p flag=enter pkg path:
	if "%flag%" == "" goto all
	yaegi extract -tag yaegi_lib %flag%
	pause
	goto input
	exit 0
:all
	for /f "tokens=1* delims=:" %%a in ('findstr /n .* yaegi_pkg_list.txt') do (
		if "%%b" == "" (pause) else (
			if not exist %%b (echo %%b && yaegi extract -tag yaegi_lib %%b) else (echo skip %%b)
		)
	)
:exit0
echo 生成完成
pause
exit 0