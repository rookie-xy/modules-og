/*
 * Copyright (C) 2017 Meng Shi
 */

package httpd

import (
      "unsafe"
    . "github.com/rookie-xy/worker/types"
    . "github.com/rookie-xy/worker/modules"
)

const (
    HTTPD_MODULE = INPUT_MODULE|0x02000000
    HTTPD_CONFIG = 0x00000002
)

var httpdModule = String{ len("httpd_module"), "httpd_module" }
var inputHttpdContext = &Context{
    httpdModule,
    nil,
    nil,
}

var httpd = String{ len("httpd"), "httpd" }
var inputHttpdCommands = []Command{

    { httpd,
      USER_CONFIG|CONFIG_BLOCK,
      httpdBlock,
      0,
      0,
      nil },

    NilCommand,
}

func httpdBlock(cycle *Cycle, _ *Command, _ *unsafe.Pointer) int {
    cycle.Configure.Block(HTTPD_MODULE, HTTPD_CONFIG)
    return Ok
}

var inputHttpdModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    unsafe.Pointer(inputHttpdContext),
    inputHttpdCommands,
    INPUT_MODULE,
    nil,
    nil,
}

func init() {
    Modules = append(Modules, &inputHttpdModule)
}
