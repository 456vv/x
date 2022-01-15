@echo off

go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://gocenter.io,https://proxy.golang.org,https://goproxy.io,https://athens.azurefd.net,direct
go mod tidy

if "%1" == "" goto input

for /f "delims=" %%I in (%1) do (yaegi extract %%I)
goto exit

:input
	set /p flag=enter pkg path:
	if "%flag%" == "" goto all
	yaegi extract %flag%
	pause
	goto input
	exit 0
:all
	for /f "tokens=1* delims=:" %%a in ('findstr /n .* yaegi_pkg_list.txt') do (
		if "%%b" == "" (pause) else (
			if not exist %%b (echo %%b && yaegi extract %%b) else (echo skip %%b)
		)
	)
	rem set fileName=../../default.go
	rem set pkgName=goplus_lib
	rem set pkgDir=github.com/456vv/x/goplus_lib
	rem 
	rem echo package %pkgName%> %fileName%
	rem echo.>> %fileName%
	rem echo import(>> %fileName%
	rem 	for /f "delims=" %%I in (yaegi_pkg_list.txt) do (echo      _ "%pkgDir%/%%I">> %fileName%)
	rem echo )>> %fileName%
	rem echo.>> %fileName%
:exit0
echo 生成完成
pause
exit 0