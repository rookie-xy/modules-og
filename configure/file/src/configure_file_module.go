/*
 * Copyright (C) 2017 Meng Shi
 */

package file

import (
    . "github.com/rookie-xy/worker/types"
    "strings"
    "fmt"
)

type fileConfigure struct {
    *Configure

     resource     string
     fileName     string
}

func NewFileConfigure(configure *Configure) *fileConfigure {
    fc := &fileConfigure{}
    fc.Configure = configure
    return fc
}

func (fc *fileConfigure) SetConfigure(configure *Configure) int {
    if configure == nil {
        return Error
    }

    fc.Configure = configure

    return Ok
}

func (fc *fileConfigure) GetConfigure() *Configure {
    return fc.Configure
}

func (fc *fileConfigure) SetResource(resource string) int {
    if resource == "" {
        return Error
    }

    fc.resource = resource

    return Ok
}

func (fc *fileConfigure) GetResource() string {
    if fc.resource == "" {
        return ""
    }

    return fc.resource
}

func (fc *fileConfigure) SetFileName(fileName string) int {
    if fileName == "" {
        return Error
    }

    fc.fileName = fileName

    return Ok
}

func (fc *fileConfigure) GetFileName() string {
    return fc.fileName
}

func (fc *fileConfigure) Set() int {
    log := fc.Log

    file := fc.File
    if file == nil {
        file = NewFile(fc.Log)
    }

    flag := false
    if file.Read() == Error {
        fmt.Println("configure read file error")
        //log.Error("configure read file error")
        flag = true
    }

    if file.Close() == Error {
        log.Warn("file close error: %d\n", 10)
        return Error
    }

    if flag {
        return Error
    }

    return Ok
}

func (fc *fileConfigure) Get() int {
    log := fc.Log
/*
    file := fc.File.Get()
    if file == nil {
        file = NewFile(fc.Log)
    }
    */

    file := NewFile(fc.Log)

    resource := fc.GetResource()
    if resource == "" {
        return Error
    }

    if file.Open(resource) == Error {
        log.Error("configure open file error")
        return Error
    }

    fc.File = file

    return Ok
}

func fileConfigureInit(cycle *Cycle) int {
    log := cycle.Log

    option := cycle.GetOption()
    if option == nil {
        return Error
    }

    item := option.GetItem("configure")
    if item == nil {
        fmt.Println("item is null")
        return Error
    }

    file := item.(string)

    fileType := file[0 : strings.Index(file, ":")]
    if fileType == "" {
        return Error
    }

    if fileType != "file" {
        return Again
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        if configure = NewConfigure(cycle); configure == nil {
            return Error
        }
    }

    fileConfigure := NewFileConfigure(configure)
    if fileConfigure == nil {
        log.Error("new file configure error")
        return Error
    }

    if fileConfigure.SetFileType(fileType) == Error {
        return Error
    }

    fileName := file[strings.LastIndex(file, "/") + 1 : ]
    if fileName == "" {
        fmt.Println("filename")
        return Error
    }

    if fileConfigure.SetFileName(fileName) == Error {
        fmt.Println("set file name")
        return Error
    }

    resource := file[strings.Index(file, "=") + 1 : ]
    if resource == "" {
        fmt.Println("resource")
        return Error
    }

    if fileConfigure.SetResource(resource) == Error {
        fmt.Println("SET RESOURCE")
        return Error
    }

    if cycle.SetConfigure(configure) == Error {
        fmt.Println("set configure")
        return Error
    }

    if configure.SetContent(fileConfigure) == Error {
        return Error
    }

    return Ok
}

func fileConfigureMain(cycle *Cycle) int {
    flag := Error
    configure := cycle.GetConfigure()
    if configure == nil {
        return flag
    }

    content := configure.GetContent()
    if content == nil {
        return flag
    }

    if content.Get() == Error {
        return flag
    }

    if content.Set() == Error {
        return flag
    }

    if flag == Error {
        fmt.Println("NNNNNNNNNNNNNNNNN")
        configure.SetNotice(Ok)
    }

    return Ok
}

var FileConfigureModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    nil,
    nil,
    SYSTEM_MODULE,
    fileConfigureInit,
    fileConfigureMain,
}

func init() {
    Modules = Load(Modules, &FileConfigureModule)
}