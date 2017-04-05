/*
 * Copyright (C) 2017 Meng Shi
 */

package file

import (
    . "github.com/rookie-xy/worker/types"
    "strings"
)

type fileConfigure struct {
    *Configure
}

func NewFileConfigure(configure *Configure) *fileConfigure {
    return &fileConfigure{ configure }
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

func (fc *fileConfigure) Set() int {
    return Ok
}

func (fc *fileConfigure) Get() int {
    return Ok
}

func initFileConfigureModule(cycle *Cycle) int {
    log := cycle.Log

    option := cycle.GetOption()
    configure := cycle.GetConfigure()
    if option == nil || configure == nil {
        return Error
    }

    item := option.GetItem("configure")
    if item == nil {
        return Error
    }

    file := item.(string)

    fileType := file[0 : strings.Index(file, ":")]
    if fileType == "" {
        return Error
    }

    if fileType != "file" {
        return Ok
    }

    if configure.SetFileType(fileType) == Error {
        return Error
    }

    fileName := file[strings.LastIndex(file, "/") + 1 : ]
    if fileName == "" {
        return Error
    }

    if configure.SetFileName(fileName) == Error {
        return Error
    }

    /*
    if configure == nil {
        configure = NewConfigure(log)
        if configure == nil {
            return Error
        }

        if cycle.SetConfigure(configure) == Error {
            return Error
        }
    }
    */

    fileConfigure := NewFileConfigure(configure)
    if fileConfigure == nil {
        log.Error("new file configure error")
        return Error
    }

    if configure.SetContent(fileConfigure) == Error {
        log.Error("set configure interface error")
        return Error
    }

    return Ok
}

var FileConfigureModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    nil,
    nil,
    SYSTEM_MODULE,
    initFileConfigureModule,
    nil,
}

func init() {
    Modules = Load(Modules, &FileConfigureModule)
}