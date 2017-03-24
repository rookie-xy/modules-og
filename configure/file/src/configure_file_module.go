/*
 * Copyright (C) 2017 Meng Shi
 */

package file

import (
    . "github.com/rookie-xy/worker/types"
)

type fileConfigure struct {
    *Configure
     Content
}

func NewFileConfigure(configure *Configure) *fileConfigure {
    return &fileConfigure{ configure, nil }
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

    cycle.Option.GetItem()

    configure := cycle.GetConfigure()
    if configure == nil {
        configure = NewConfigure(log)
        if configure == nil {
            return Error
        }

        if cycle.SetConfigure(configure) == Error {
            return Error
        }
    }

    fileConfigure := NewFileConfigure(configure)
    if fileConfigure == nil {
        log.Error("new file configure error")
        return Error
    }

    if configure.Set(fileConfigure) == Error {
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
    Modules = append(Modules, &FileConfigureModule)
}