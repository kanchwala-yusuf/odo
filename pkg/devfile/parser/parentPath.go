package parser

import (
	"encoding/json"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// GetParentPath returns parentPath stored in devfile context
func (d *DevfileCtx) GetParentPath() string {
	return d.parentPath
}

// SetDevfileAPIVersion returns the devfile APIVersion
func (d *DevfileCtx) SetDevfileParentPath() error {

	// Unmarshal JSON into map
	var r map[string]interface{}
	err := json.Unmarshal(d.rawContent, &r)
	if err != nil {
		return errors.Wrapf(err, "failed to decode devfile json")
	}

	// check if "parent" is present in devfile
	_, ok := r["parent"]
	if !ok {
		glog.V(4).Infof("'parent' not present in devfile")
		d.parentPath = ""
		return nil
	}

	// Get "parentPath" value from the map
	parent := r["parent"].(map[string]interface{})
	parentPath, ok := parent["uri"]
	if !ok {
		glog.V(4).Infof("'parent.uri' not present in devfile")
		d.parentPath = ""
		return nil
	}

	// Successful
	d.parentPath = parentPath.(string)
	glog.V(4).Infof("devfile parentPath: '%s'", d.parentPath)
	return nil
}
