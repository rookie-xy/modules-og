/*
 * Copyright (C) 2017 Meng Shi
 */

package zookeeper

import (
    . "github.com/rookie-xy/worker/types"
    "fmt"
    "strings"
)

type zookeeperConfigure struct {
    *Configure
}

func NewZookeeperConfigure(configure *Configure) *zookeeperConfigure {
    return &zookeeperConfigure{ configure }
}

func (zkc *zookeeperConfigure) SetConfigure(configure *Configure) int {
    if configure == nil {
        return Error
    }

    zkc.Configure = configure

    return Ok
}

func (zkc *zookeeperConfigure) GetConfigure() *Configure {
    return zkc.Configure
}

func (zkc *zookeeperConfigure) Set() int {
    return Ok
}

func (zkc *zookeeperConfigure) Get() int {
    return Ok
}

func zookeeperConfigureInit(cycle *Cycle) int {
    //log := cycle.Log

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

    if fileType != "zookeeper" {
        return Ignore
    }

    return Ok
}

func zookeeperConfigureMain(cycle *Cycle) int {
    return Ok
}

var ZookeeperConfigureModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    nil,
    nil,
    SYSTEM_MODULE,
    zookeeperConfigureInit,
    zookeeperConfigureMain,
}

func init() {
    Modules = Load(Modules, &ZookeeperConfigureModule)
}