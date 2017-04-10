/*
 * Copyright (C) 2017 Meng Shi
 */

package file

import (
      "strings"
      "github.com/fsnotify/fsnotify"

    . "github.com/rookie-xy/worker/types"

    "unsafe"
    "log"
)

type fileConfigure struct {
    *Configure

     watcher     *fsnotify.Watcher

     resource     string
     fileName     string
}

func NewFileConfigure(configure *Configure) *fileConfigure {
    fc := &fileConfigure{}
    fc.Configure = configure
    return fc
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

func (fc *fileConfigure) SetResource(resource string) int {
    if resource == "" {
        return Error
    }

    fc.resource = resource

    return Ok
}

func (fc *fileConfigure) GetResource() string {
    if fc.resource == "" {
        return ""
    }

    return fc.resource
}

func (fc *fileConfigure) SetFileName(fileName string) int {
    if fileName == "" {
        return Error
    }

    fc.fileName = fileName

    return Ok
}

func (fc *fileConfigure) GetFileName() string {
    return fc.fileName
}

func (fc *fileConfigure) Set() int {
    if fc.AbstractFile == nil {
        fc.Error("file configure set error")
        return Error
    }

    flag := false
    if fc.Reader() == Error {
        fc.Error("configure read file error")
        flag = true
    }

    if fc.Closer() == Error {
        fc.Warn("file close error: %d\n", 10)
        return Error
    }

    if flag {
        return Error
    }

    return Ok
}

func (fc *fileConfigure) Get() int {
    if fc.AbstractFile == nil {
        fc.AbstractFile = NewAbstractFile(fc.Log)
    }

    resource := fc.GetResource()
    if resource == "" {
        return Error
    }

    if fc.Open(resource) == Error {
        fc.Error("configure open file error")
        return Error
    }

    return Ok
}

func (fc *fileConfigure) GetType() unsafe.Pointer {
    return unsafe.Pointer(fc)
}

func fileConfigureInit(cycle *Cycle) int {
    option := cycle.GetOption()
    if option == nil {
        return Error
    }

    item := option.GetItem("configure")
    if item == nil {
        return Error
    }

    file := item.(string)

    fileType := file[0 : strings.Index(file, ":")]
    if fileType == "" {
        return Error
    }

    if fileType != "file" {
        return Ignore
    }

    configure := cycle.GetConfigure()
    if configure == nil {
        if configure = NewConfigure(cycle); configure == nil {
            return Error
        }
    }

    fileConfigure := NewFileConfigure(configure)
    if fileConfigure == nil {
        return Error
    }

    if watcher, error := fsnotify.NewWatcher(); error != nil {
        //fileConfigure.Error(error)
        return Error
    } else {
        fileConfigure.watcher = watcher
    }

    if fileConfigure.SetName(fileType) == Error {
        return Error
    }

    fileName := file[strings.LastIndex(file, "/") + 1 : ]
    if fileName == "" {
        return Error
    }

    if fileConfigure.SetFileName(fileName) == Error {
        return Error
    }

    resource := file[strings.Index(file, "=") + 1 : ]
    if resource == "" {
        return Error
    }

    if fileConfigure.SetResource(resource) == Error {
        return Error
    }

    if error := fileConfigure.watcher.Add(resource); error != nil {
        return Error
    }

    if cycle.SetConfigure(configure) == Error {
        return Error
    }

    if configure.SetHandle(fileConfigure) == Error {
        return Error
    }

    return Ok
}

func fileConfigureMain(cycle *Cycle) int {
    flag := Error
    configure := cycle.GetConfigure()
    if configure == nil {
        return flag
    }

    content := configure.GetHandle()
    if content == nil {
        return flag
    }

    if content.Get() == Error {
        return flag
    }

    if content.Set() == Error {
        return flag
    }

    if flag == Error {
        configure.SetNotice(Ok)
    }

    fc := (*fileConfigure)(unsafe.Pointer(content.GetType()))
    if fc == nil {
        return Error
    }

    defer fc.watcher.Close()

    for {
        select {

        case event := <-fc.watcher.Events:
            log.Println("event:", event)
            if event.Op & fsnotify.Write == fsnotify.Write {
                log.Println("modified file:", event.Name)
            }

        case err := <-fc.watcher.Errors:
            log.Println("error:", err)
        }
    }

    return Ok
}

var FileConfigureModule = Module{
    MODULE_V1,
    CONTEXT_V1,
    nil,
    nil,
    SYSTEM_MODULE,
    fileConfigureInit,
    fileConfigureMain,
}

func init() {
    Modules = Load(Modules, &FileConfigureModule)
}