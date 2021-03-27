@echo off

go env -w GOPROXY=https://goproxy.cn,https://mirrors.aliyun.com/goproxy/,https://gocenter.io,https://proxy.golang.org,https://goproxy.io,https://athens.azurefd.net,direct
go mod tidy

if "%1" == "" goto input

for /f "delims=" %%I in (%1) do (gop export -outdir . %%I)
goto exit

:input
	set /p flag=enter pkg path:
	if "%flag%" == "" goto all
	gop export -outdir . %flag%
	pause
	goto input
	exit 0
:all
	for /f "tokens=1* delims=:" %%a in ('findstr /n .* gop_gen_lib_pkg_list.txt') do (
		if "%%b" == "" (pause) else (
			if not exist %%b (gop export -outdir . %%b) else (echo skip %%b)
		)
	)
	set fileName=default.go
	set pkgName=goplus_lib
	set pkgDir=github.com/456vv/x/goplus_lib
	
	echo package %pkgName%> %fileName%
	echo.>> %fileName%
	echo import(>> %fileName%
		for /f "delims=" %%I in (gop_gen_lib_pkg_list.txt) do (echo      _ "%pkgDir%/%%I">> %fileName%)
	echo )>> %fileName%
	echo.>> %fileName%
:exit0
echo 生成完成
pause
exit 0