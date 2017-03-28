/*
 * Copyright (C) 2017 Meng Shi
 */

package multiline

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    . "github.com/rookie-xy/worker/modules"
)

const (
    MULTILINE_MODULE = 0x0007
    MULTILINE_CONFIG = 0x00070000
)

var multilineModule = String{ len("multiline_module"), "multiline_module" }
var codecMultilineContext = &Context{
    multilineModule,
    nil,
    nil,
}

var	multiline = String{ len("multiline"), "multiline" }
var codecMultilineCommands = []Command{

    { multiline,
      USER_CONFIG|CONFIG_BLOCK,
      multilineBlock,
      0,
      0,
      nil },

    NilCommand,
}

func multilineBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    cycle.Configure.Block(CODEC_MODULE|MULTILINE_MODULE, MULTILINE_CONFIG)
    return Ok
}

var codecMultilineModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(codecMultilineContext),
    codecMultilineCommands,
    CODEC_MODULE,
    nil,
    nil,
}

func init() {
    Modules = Load(Modules, &codecMultilineModule)
}