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
    MULTILINE_MODULE = CODEC_MODULE|0x01000000
    MULTILINE_CONFIG = USER_CONFIG|CONFIG_MAP
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
      MULTILINE_CONFIG,
      multilineBlock,
      0,
      0,
      nil },

    NilCommand,
}

func multilineBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    cycle.Configure.Block(MULTILINE_MODULE, MULTILINE_CONFIG|CONFIG_VALUE)
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