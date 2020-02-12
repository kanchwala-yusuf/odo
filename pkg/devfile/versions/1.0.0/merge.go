package version100

import (
	"reflect"

	"github.com/openshift/odo/pkg/devfile/versions/common"
)

// IsFieldPresentInStruct checks if the provided arg interface is either a
// struct or a pointer to a struct and has the defined member field.
// Returns:
// true: if FieldName exists and is accessible with reflect
// false: if FieldName is not present
func IsFieldPresentInStruct(Iface interface{}, FieldName string) bool {
	ValueIface := reflect.ValueOf(Iface)

	// Check if the passed interface is a pointer
	if ValueIface.Type().Kind() != reflect.Ptr {
		// Create a new type of Iface's Type, so we have a pointer to work with
		ValueIface = reflect.New(reflect.TypeOf(Iface))
	}

	// 'dereference' with Elem() and get the field by name
	Field := ValueIface.Elem().FieldByName(FieldName)
	if !Field.IsValid() {
		return false
	}
	return true
}

// MergeDevfile combines parent and local devfiles
func (local *Devfile100) MergeDevfiles(p interface{}) error {

	// devfile fields
	fields := []string{"Metadata", "Projects", "Components", "Commands"}

	// new devfile data object
	parent := p.(*Devfile100)

	// merge all fields
	for _, f := range fields {
		switch f {
		case "Metadata":
			if !(IsFieldPresentInStruct(local, f) && !reflect.DeepEqual(local.Metadata, common.DevfileMetadata{})) {
				local.Metadata = parent.Metadata
			}
		case "Projects":
			// If parent devfile does not have "Projects" field, then nothing needs
			// to be done
			if IsFieldPresentInStruct(parent, f) {
				local.mergeProjects(parent.Projects)
			}
		case "Components":
			// If parent devfile does not have "Components" field, then nothing needs
			// to be done
			if IsFieldPresentInStruct(parent, f) {
				local.mergeComponents(parent.Components)
			}
		case "Commands":
			// If parent devfile does not have "Commands" field, then nothing needs
			// to be done
			if IsFieldPresentInStruct(parent, f) {
				local.mergeCommands(parent.Commands)
			}
		}
	}
	return nil
}

func (local *Devfile100) mergeProjects(parentProjects []common.DevfileProject) error {

	// if list of projects is empty in parent devfile, then nothing needs to be
	// overriden
	if len(parentProjects) < 1 {
		return nil
	}

	// if list of projects is empty in local devfile, then merge parent projects
	if len(local.Projects) < 1 {
		local.Projects = append(local.Projects, parentProjects...)
		return nil
	}

	// make a map of project names and project type
	localProjectsMap := make(map[string]bool)
	for _, project := range local.Projects {
		localProjectsMap[project.Name] = true
	}

	// if parent has project which are not in local devfile
	for _, project := range parentProjects {
		if _, ok := localProjectsMap[project.Name]; !ok {
			local.Projects = append(local.Projects, project)
		}
	}

	// successfull
	return nil
}

func (local *Devfile100) mergeCommands(parentCommands []common.DevfileCommand) error {

	// if list of parent commands is empty, then nothing needs to be done
	if len(parentCommands) < 1 {
		return nil
	}

	// if list of local commands is empty, then merge parent commands
	if len(local.Commands) < 1 {
		local.Commands = append(local.Commands, parentCommands...)
		return nil
	}

	// make a map of command names and command type
	localCommandsMap := make(map[string]bool)
	for _, command := range local.Commands {
		localCommandsMap[command.Name] = true
	}

	// if parent has command which are not in local devfile
	for _, command := range parentCommands {
		if _, ok := localCommandsMap[command.Name]; !ok {
			local.Commands = append(local.Commands, command)
		}
	}

	// successfull
	return nil
}

func (local *Devfile100) mergeComponents(parentComponents []common.DevfileComponent) error {

	// if list of parent components is empty, then nothing needs to be done
	if len(parentComponents) < 1 {
		return nil
	}

	// if list of local components is empty, then merge parent components
	if len(local.Components) < 1 {
		local.Components = append(local.Components, parentComponents...)
		return nil
	}

	// make a map of component names and component type
	localComponentsMap := make(map[string]bool)
	for _, component := range local.Components {

		if component.Type == common.DevfileComponentTypeCheEditor || component.Type == common.DevfileComponentTypeChePlugin {
			localComponentsMap[string(component.Type)] = true
		} else if component.Alias != nil {
			localComponentsMap[*component.Alias] = true
		}
	}

	// if parent has component which are not in local devfile
	for _, component := range parentComponents {
		if component.Type == common.DevfileComponentTypeCheEditor || component.Type == common.DevfileComponentTypeChePlugin {
			if _, ok := localComponentsMap[string(component.Type)]; !ok {
				local.Components = append(local.Components, component)
			}
		} else if component.Id != nil {
			if _, ok := localComponentsMap[*component.Id]; !ok {
				local.Components = append(local.Components, component)
			}
		}
	}

	// successfull
	return nil
}
