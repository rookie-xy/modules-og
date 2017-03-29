/*
 * Copyright (C) 2017 Meng Shi
 */

package mlog

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    . "github.com/rookie-xy/worker/modules"
)

const (
    MLOG_MODULE = LOG_MODULE|0x01000000
    MLOG_CONFIG = 0x00020000
)

var mlogModule = String{ len("mlog_module"), "mlog_module" }
var logMlogContext = &Context{
    mlogModule,
    nil,
    nil,
}

var mlog = String{ len("mlog"), "mlog" }
var logMlogCommands = []Command{

    { mlog,
      USER_CONFIG|CONFIG_BLOCK,
      mlogBlock,
      0,
      0,
      nil },

    NilCommand,
}

func mlogBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    cycle.Configure.Block(MLOG_MODULE, MLOG_CONFIG)
    return Ok
}

var logMlogModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(logMlogContext),
    logMlogCommands,
    LOG_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &logMlogModule)
}