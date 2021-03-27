// Package fs provide Go+ "io/fs" package, as "io/fs" package in Go.
package fs

import (
	fs "io/fs"
	reflect "reflect"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execiDirEntryInfo(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(fs.DirEntry).Info()
	p.Ret(1, ret0, ret1)
}

func execiDirEntryIsDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.DirEntry).IsDir()
	p.Ret(1, ret0)
}

func execiDirEntryName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.DirEntry).Name()
	p.Ret(1, ret0)
}

func execiDirEntryType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.DirEntry).Type()
	p.Ret(1, ret0)
}

func execiFSOpen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.FS).Open(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiFileClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.File).Close()
	p.Ret(1, ret0)
}

func execiFileRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.File).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiFileStat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(fs.File).Stat()
	p.Ret(1, ret0, ret1)
}

func execiFileInfoIsDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileInfo).IsDir()
	p.Ret(1, ret0)
}

func execiFileInfoModTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileInfo).ModTime()
	p.Ret(1, ret0)
}

func execiFileInfoMode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileInfo).Mode()
	p.Ret(1, ret0)
}

func execiFileInfoName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileInfo).Name()
	p.Ret(1, ret0)
}

func execiFileInfoSize(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileInfo).Size()
	p.Ret(1, ret0)
}

func execiFileInfoSys(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileInfo).Sys()
	p.Ret(1, ret0)
}

func execmFileModeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileMode).String()
	p.Ret(1, ret0)
}

func execmFileModeIsDir(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileMode).IsDir()
	p.Ret(1, ret0)
}

func execmFileModeIsRegular(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileMode).IsRegular()
	p.Ret(1, ret0)
}

func execmFileModePerm(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileMode).Perm()
	p.Ret(1, ret0)
}

func execmFileModeType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(fs.FileMode).Type()
	p.Ret(1, ret0)
}

func toType0(v interface{}) fs.FS {
	if v == nil {
		return nil
	}
	return v.(fs.FS)
}

func execGlob(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := fs.Glob(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiGlobFSGlob(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.GlobFS).Glob(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmPathErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*fs.PathError).Error()
	p.Ret(1, ret0)
}

func execmPathErrorUnwrap(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*fs.PathError).Unwrap()
	p.Ret(1, ret0)
}

func execmPathErrorTimeout(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*fs.PathError).Timeout()
	p.Ret(1, ret0)
}

func execReadDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := fs.ReadDir(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiReadDirFSReadDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.ReadDirFS).ReadDir(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiReadDirFileReadDir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.ReadDirFile).ReadDir(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execReadFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := fs.ReadFile(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiReadFileFSReadFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.ReadFileFS).ReadFile(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execStat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := fs.Stat(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiStatFSStat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.StatFS).Stat(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execSub(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := fs.Sub(toType0(args[0]), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiSubFSSub(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(fs.SubFS).Sub(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execValidPath(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := fs.ValidPath(args[0].(string))
	p.Ret(1, ret0)
}

func execWalkDir(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := fs.WalkDir(toType0(args[0]), args[1].(string), args[2].(fs.WalkDirFunc))
	p.Ret(3, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("io/fs")

func init() {
	I.RegisterFuncs(
		I.Func("(DirEntry).Info", (fs.DirEntry).Info, execiDirEntryInfo),
		I.Func("(DirEntry).IsDir", (fs.DirEntry).IsDir, execiDirEntryIsDir),
		I.Func("(DirEntry).Name", (fs.DirEntry).Name, execiDirEntryName),
		I.Func("(DirEntry).Type", (fs.DirEntry).Type, execiDirEntryType),
		I.Func("(FS).Open", (fs.FS).Open, execiFSOpen),
		I.Func("(File).Close", (fs.File).Close, execiFileClose),
		I.Func("(File).Read", (fs.File).Read, execiFileRead),
		I.Func("(File).Stat", (fs.File).Stat, execiFileStat),
		I.Func("(FileInfo).IsDir", (fs.FileInfo).IsDir, execiFileInfoIsDir),
		I.Func("(FileInfo).ModTime", (fs.FileInfo).ModTime, execiFileInfoModTime),
		I.Func("(FileInfo).Mode", (fs.FileInfo).Mode, execiFileInfoMode),
		I.Func("(FileInfo).Name", (fs.FileInfo).Name, execiFileInfoName),
		I.Func("(FileInfo).Size", (fs.FileInfo).Size, execiFileInfoSize),
		I.Func("(FileInfo).Sys", (fs.FileInfo).Sys, execiFileInfoSys),
		I.Func("(FileMode).String", (fs.FileMode).String, execmFileModeString),
		I.Func("(FileMode).IsDir", (fs.FileMode).IsDir, execmFileModeIsDir),
		I.Func("(FileMode).IsRegular", (fs.FileMode).IsRegular, execmFileModeIsRegular),
		I.Func("(FileMode).Perm", (fs.FileMode).Perm, execmFileModePerm),
		I.Func("(FileMode).Type", (fs.FileMode).Type, execmFileModeType),
		I.Func("Glob", fs.Glob, execGlob),
		I.Func("(GlobFS).Glob", (fs.GlobFS).Glob, execiGlobFSGlob),
		I.Func("(*PathError).Error", (*fs.PathError).Error, execmPathErrorError),
		I.Func("(*PathError).Unwrap", (*fs.PathError).Unwrap, execmPathErrorUnwrap),
		I.Func("(*PathError).Timeout", (*fs.PathError).Timeout, execmPathErrorTimeout),
		I.Func("ReadDir", fs.ReadDir, execReadDir),
		I.Func("(ReadDirFS).ReadDir", (fs.ReadDirFS).ReadDir, execiReadDirFSReadDir),
		I.Func("(ReadDirFile).ReadDir", (fs.ReadDirFile).ReadDir, execiReadDirFileReadDir),
		I.Func("ReadFile", fs.ReadFile, execReadFile),
		I.Func("(ReadFileFS).ReadFile", (fs.ReadFileFS).ReadFile, execiReadFileFSReadFile),
		I.Func("Stat", fs.Stat, execStat),
		I.Func("(StatFS).Stat", (fs.StatFS).Stat, execiStatFSStat),
		I.Func("Sub", fs.Sub, execSub),
		I.Func("(SubFS).Sub", (fs.SubFS).Sub, execiSubFSSub),
		I.Func("ValidPath", fs.ValidPath, execValidPath),
		I.Func("WalkDir", fs.WalkDir, execWalkDir),
	)
	I.RegisterVars(
		I.Var("ErrClosed", &fs.ErrClosed),
		I.Var("ErrExist", &fs.ErrExist),
		I.Var("ErrInvalid", &fs.ErrInvalid),
		I.Var("ErrNotExist", &fs.ErrNotExist),
		I.Var("ErrPermission", &fs.ErrPermission),
		I.Var("SkipDir", &fs.SkipDir),
	)
	I.RegisterTypes(
		I.Type("DirEntry", reflect.TypeOf((*fs.DirEntry)(nil)).Elem()),
		I.Type("FS", reflect.TypeOf((*fs.FS)(nil)).Elem()),
		I.Type("File", reflect.TypeOf((*fs.File)(nil)).Elem()),
		I.Type("FileInfo", reflect.TypeOf((*fs.FileInfo)(nil)).Elem()),
		I.Type("FileMode", reflect.TypeOf((*fs.FileMode)(nil)).Elem()),
		I.Type("GlobFS", reflect.TypeOf((*fs.GlobFS)(nil)).Elem()),
		I.Type("PathError", reflect.TypeOf((*fs.PathError)(nil)).Elem()),
		I.Type("ReadDirFS", reflect.TypeOf((*fs.ReadDirFS)(nil)).Elem()),
		I.Type("ReadDirFile", reflect.TypeOf((*fs.ReadDirFile)(nil)).Elem()),
		I.Type("ReadFileFS", reflect.TypeOf((*fs.ReadFileFS)(nil)).Elem()),
		I.Type("StatFS", reflect.TypeOf((*fs.StatFS)(nil)).Elem()),
		I.Type("SubFS", reflect.TypeOf((*fs.SubFS)(nil)).Elem()),
		I.Type("WalkDirFunc", reflect.TypeOf((*fs.WalkDirFunc)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ModeAppend", qspec.Uint32, fs.ModeAppend),
		I.Const("ModeCharDevice", qspec.Uint32, fs.ModeCharDevice),
		I.Const("ModeDevice", qspec.Uint32, fs.ModeDevice),
		I.Const("ModeDir", qspec.Uint32, fs.ModeDir),
		I.Const("ModeExclusive", qspec.Uint32, fs.ModeExclusive),
		I.Const("ModeIrregular", qspec.Uint32, fs.ModeIrregular),
		I.Const("ModeNamedPipe", qspec.Uint32, fs.ModeNamedPipe),
		I.Const("ModePerm", qspec.Uint32, fs.ModePerm),
		I.Const("ModeSetgid", qspec.Uint32, fs.ModeSetgid),
		I.Const("ModeSetuid", qspec.Uint32, fs.ModeSetuid),
		I.Const("ModeSocket", qspec.Uint32, fs.ModeSocket),
		I.Const("ModeSticky", qspec.Uint32, fs.ModeSticky),
		I.Const("ModeSymlink", qspec.Uint32, fs.ModeSymlink),
		I.Const("ModeTemporary", qspec.Uint32, fs.ModeTemporary),
		I.Const("ModeType", qspec.Uint32, fs.ModeType),
	)
}
