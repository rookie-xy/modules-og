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
    GROK_MODULE = 0x0007
    GROK_CONFIG = 0x00070000
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
      USER_CONFIG|CONFIG_BLOCK,
      grokBlock,
      0,
      0,
      nil },

    NilCommand,
}

func grokBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    cycle.Configure.Block(FILTER_MODULE|GROK_MODULE, GROK_CONFIG)
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
    Modules = append(Modules, &filterGrokModule)
}