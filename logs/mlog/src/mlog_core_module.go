/*
 * Copyright (C) 2017 Meng Shi
 */

package mlog

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
)

type MlogCore struct {
    *Log

     level string
     path  string
     log   LogIf
}

func NewMlogCore() *MlogCore {
    return &MlogCore{}
}

var mlogCore = String{ len("mlog_core"), "mlog_core" }
var coreMlogContext = &Context{
    mlogCore,
    coreMlogContextCreate,
    coreMlogContextInit,
}

func coreMlogContextCreate(cycle *Cycle) unsafe.Pointer {
    mlogCore := NewMlogCore()
    if mlogCore == nil {
        return nil
    }

    mlogCore.level = "debug"
    mlogCore.path = "/data/logs/service"

    return unsafe.Pointer(mlogCore)
}

func coreMlogContextInit(cycle *Cycle, context *unsafe.Pointer) string {
    log := cycle.GetLog()
    this := (*MlogCore)(unsafe.Pointer(uintptr(*context)))
    if this == nil {
        log.Error("coreMlogContextInit error")
        return "0"
    }

    level, path := this.level, this.path
    if level == "" || path == "" {
        return "0"
    }

    for i, v := range Level {
        if v == level {
            if log.SetLevel(i) == Error {
                return "error"
            }
        }
    }

    if log.SetPath(path) == Error {
        return "error"
    }

    return "0"
}

var (
    level = String{ len("level"), "level" }
    path = String{ len("path"), "path" }
    coreMlog MlogCore
)

var coreMlogCommands = []Command{

    { level,
      MLOG_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(coreMlog.level),
      nil },

    { path,
      MLOG_CONFIG,
      SetString,
      0,
      unsafe.Offsetof(coreMlog.path),
      nil },

    NilCommand,
}

var coreMlogModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(coreMlogContext),
    coreMlogCommands,
    MLOG_MODULE,
    coreMlogInit,
    coreMlogMain,
}

func coreMlogInit(cycle *Cycle) int {
    //log := cycle.Log
    //fmt.Println(log.GetPath())
    return Ok
}

func coreMlogMain(cycle *Cycle) int {
    return Ok
}

func init() {
    Modules = append(Modules, &coreMlogModule)
}