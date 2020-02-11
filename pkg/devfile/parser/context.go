package parser

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift/odo/pkg/testingutil/filesystem"
	"github.com/openshift/odo/pkg/util"
)

// DevfileCtx stores context info regarding devfile
type DevfileCtx struct {

	// absolute path of devfile
	absPath string

	// devfile ApiVersion
	apiVersion string

	// filesystem for devfile
	Fs filesystem.Filesystem

	// devfile json schema
	jsonSchema string

	// parent devfile path
	parentPath string

	// raw content of the devfile
	rawContent []byte

	// relative path of devfile
	relPath string
}

// NewDevfileCtx returns a new DevfileCtx type object
func NewDevfileCtx(path string) DevfileCtx {
	return DevfileCtx{
		relPath: path,
		Fs:      filesystem.DefaultFs{},
	}
}

// Populate fills the DevfileCtx struct with relevant context info
func (d *DevfileCtx) Populate() (err error) {

	// Get devfile absolute path
	if d.absPath, err = util.GetAbsPath(d.relPath); err != nil {
		return err
	}
	glog.V(4).Infof("absolute devfile path: '%s'", d.absPath)

	// Read and save devfile content
	if err := d.SetDevfileContent(); err != nil {
		return err
	}

	// Get devfile APIVersion
	if err := d.SetDevfileAPIVersion(); err != nil {
		return err
	}

	// Check if the apiVersion is supported
	if !d.IsApiVersionSupported() {
		return fmt.Errorf("devfile apiVersion '%s' not supported in odo", d.apiVersion)
	}
	glog.V(4).Infof("devfile apiVersion '%s' is supported in odo", d.apiVersion)

	// Get devfile parent
	if err := d.SetDevfileParentPath(); err != nil {
		return err
	}

	// Read and save devfile JSON schema for provided apiVersion
	if err := d.SetDevfileJSONSchema(); err != nil {
		return err
	}

	// Successful
	return nil
}

// Validate func validates devfile JSON schema for the given apiVersion
func (d *DevfileCtx) Validate() error {

	// Validate devfile
	return d.ValidateDevfileSchema()
}
