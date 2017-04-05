/*
 * Copyright (C) 2017 Meng Shi
 */

package grok

import (
      "unsafe"

    . "github.com/rookie-xy/worker/types"
    . "github.com/rookie-xy/worker/modules"
)

const (
    GROK_MODULE = FILTER_MODULE|0x01000000
    GROK_CONFIG = USER_CONFIG|CONFIG_ARRAY
)

var grokModule = String{ len("grok_module"), "grok_module" }
var filterGrokContext = &Context{
    grokModule,
    nil,
    nil,
}

var	grok = String{ len("grok"), "grok" }
var filterGrokCommands = []Command{

    { grok,
      GROK_CONFIG,
      grokBlock,
      0,
      0,
      nil },

    NilCommand,
}

func grokBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    cycle.Configure.Block(GROK_MODULE, GROK_CONFIG|CONFIG_VALUE)
    return Ok
}

var filterGrokModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(filterGrokContext),
    filterGrokCommands,
    FILTER_MODULE,
    nil,
    nil,
}

func init() {
    Modules = Load(Modules, &filterGrokModule)
}