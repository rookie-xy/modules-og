/*
 * Copyright (C) 2017 Meng Shi
 */

package zookeeper

import (
    . "github.com/rookie-xy/worker/types"
)

type zookeeperConfigure struct {
    *Configure
     Content
}

func NewZookeeperConfigure(configure *Configure) *zookeeperConfigure {
    return &zookeeperConfigure{ configure, nil }
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

func initZookeeperConfigureModule(cycle *Cycle) int {
    log := cycle.Log

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

    zookeeperConfigure := NewZookeeperConfigure(configure)
    if zookeeperConfigure == nil {
        log.Error("new zookeeper configure error")
        return Error
    }

    if configure.Set(zookeeperConfigure) == Error {
        log.Error("set configure interface error")
        return Error
    }

    return Ok
}

var ZookeeperConfigureModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    nil,
    nil,
    SYSTEM_MODULE,
    initZookeeperConfigureModule,
    nil,
}

func init() {
    Modules = append(Modules, &ZookeeperConfigureModule)
}