package devfile

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/openshift/odo/pkg/devfile/parser"
	"github.com/openshift/odo/pkg/devfile/versions"
	"github.com/pkg/errors"
)

// Parse func parses and validates the devfile integrity.
// Creates devfile context and runtime objects
func Parse(path string) (d DevfileObj, err error) {

	// NewDevfileCtx
	d.Ctx = parser.NewDevfileCtx(path)

	// Fill the fields of DevfileCtx struct
	err = d.Ctx.Populate()
	if err != nil {
		return d, err
	}

	// Create a new devfile data object
	d.Data, err = versions.NewDevfileData(d.Ctx.GetApiVersion())
	if err != nil {
		return d, err
	}

	// Unmarshal devfile content into devfile struct
	err = json.Unmarshal(d.Ctx.GetDevfileContent(), &d.Data)
	if err != nil {
		return d, errors.Wrapf(err, "failed to decode devfile content")
	}

	// if parent devfile is present, fetch and use it as base devfile
	if d.Ctx.GetParentPath() != "" {
		glog.V(4).Infof("processing parent devfile")

		// parse parent devfile
		p, err := Parse(d.Ctx.GetParentPath())
		if err != nil {
			return d, errors.Wrapf(err, "failed to parse parent devfile")
		}

		// parent and local devfile apiVersions cannot be different
		if p.Ctx.GetApiVersion() != d.Ctx.GetApiVersion() {
			return d, errors.Errorf("parent and local devfiles version mismatch")
		}

		// merge parent and local devfiles
		err = d.Data.MergeDevfiles(p.Data)
		if err != nil {
			return d, errors.Wrapf(err, "failed to merge parent and local devfiles")
		}
	}

	// Validate devfile context
	err = d.Ctx.Validate()
	if err != nil {
		return d, err
	}

	// Validate devfile data
	err = d.Data.Validate()
	if err != nil {
		return d, err
	}

	// Successful
	return d, nil
}
