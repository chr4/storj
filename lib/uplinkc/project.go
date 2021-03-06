// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

// #include "uplink_definitions.h"
import "C"
import (
	libuplink "storj.io/storj/lib/uplink"
)

// Project is a scoped uplink.Project
type Project struct {
	scope
	*libuplink.Project
}

//export open_project
// open_project opens project using uplink
func open_project(uplinkHandle C.UplinkRef, satelliteAddr *C.char, apikeyHandle C.APIKeyRef, cerr **C.char) C.ProjectRef {
	uplink, ok := universe.Get(uplinkHandle._handle).(*Uplink)
	if !ok {
		*cerr = C.CString("invalid uplink")
		return C.ProjectRef{}
	}

	apikey, ok := universe.Get(apikeyHandle._handle).(libuplink.APIKey)
	if !ok {
		*cerr = C.CString("invalid apikey")
		return C.ProjectRef{}
	}

	scope := uplink.scope.child()

	project, err := uplink.OpenProject(scope.ctx, C.GoString(satelliteAddr), apikey)
	if err != nil {
		*cerr = C.CString(err.Error())
		return C.ProjectRef{}
	}

	return C.ProjectRef{universe.Add(&Project{scope, project})}
}

//export close_project
// close_project closes the project.
func close_project(projectHandle C.ProjectRef, cerr **C.char) {
	project, ok := universe.Get(projectHandle._handle).(*Project)
	if !ok {
		*cerr = C.CString("invalid project")
		return
	}
	universe.Del(projectHandle._handle)
	defer project.cancel()

	if err := project.Close(); err != nil {
		*cerr = C.CString(err.Error())
		return
	}
}
