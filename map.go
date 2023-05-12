package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type Action string
type ResourceType string

const (
	ResourceTypeFile     ResourceType = "file"
	ResourceTypeResource ResourceType = "resource"
	ResourceTypeModule   ResourceType = "module"
	DefaultFileName      string       = "unknown file"
)

const (
	// ActionNoop denotes a no-op operation.
	ActionNoop Action = "no-op"

	// ActionCreate denotes a create operation.
	ActionCreate Action = "create"

	// ActionRead denotes a read operation.
	ActionRead Action = "read"

	// ActionUpdate denotes an update operation.
	ActionUpdate Action = "update"

	// ActionDelete denotes a delete operation.
	ActionDelete Action = "delete"

	// ActionReplace denotes a replace operation.
	ActionReplace Action = "replace"
)

// Map represents the root module
type Map struct {
	Path              string                                   `json:"path"`
	RequiredCore      []string                                 `json:"required_core,omitempty"`
	RequiredProviders map[string]*tfconfig.ProviderRequirement `json:"required_providers,omitempty"`
	// ProviderConfigs   map[string]*tfconfig.ProviderConfig      `json:"provider_configs,omitempty"`
	Root map[string]*Resource `json:"root,omitempty"`
}

// Resource is a modified tfconfig.Resource
type Resource struct {
	Type ResourceType `json:"type"`
	Name string       `json:"name"`
	Line *int         `json:"line,omitempty"`

	Children map[string]*Resource `json:"children,omitempty"`

	// Resource
	ChangeAction Action `json:"change_action,omitempty"`
	// Variable and Output
	Required  *bool `json:"required,omitempty"`
	Sensitive bool  `json:"sensitive,omitempty"`
	// Provider and Data
	Provider     string `json:"provider,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
	// ModuleCall
	Source  string `json:"source,omitempty"`
	Version string `json:"version,omitempty"`
}

// ModuleCall is a modified tfconfig.ModuleCall
type ModuleCall struct {
	Name    string `json:"name"`
	Source  string `json:"source"`
	Version string `json:"version,omitempty"`
	Line    int    `json:"line,omitempty"`
}

func (r *rover) GenerateModuleMap(parent *Resource, parentModule string) {

	childIndex := regexp.MustCompile(`\[[^[\]]*\]$`)
	matchBrackets := regexp.MustCompile(`\[[^\[\]]*\]`)

	states := r.RSO.States
	configs := r.RSO.Configs

	prefix := parentModule
	if parentModule != "" {
		prefix = fmt.Sprintf("%s.", prefix)
	}

	parentConfig := matchBrackets.ReplaceAllString(parentModule, "")

	for id, rs := range states[parentModule].Children {

		configId := matchBrackets.ReplaceAllString(id, "")
		configured := configs[parentConfig] != nil && configs[parentConfig].Module != nil && configs[configId] != nil // If there is configuration for filenames, lines, etc.

		re := &Resource{
			Type:     rs.Type,
			Children: map[string]*Resource{},
		}

		if states[id].Change.Actions != nil {

			re.ChangeAction = Action(string(states[id].Change.Actions[0]))
			if len(states[id].Change.Actions) > 1 {
				re.ChangeAction = ActionReplace
			}
		}

		if rs.Type == ResourceTypeResource {
			re.ResourceType = configs[configId].ResourceConfig.Type
			re.Name = configs[configId].ResourceConfig.Name

			for crName, cr := range states[id].Children {

				if re.Children == nil {
					re.Children = make(map[string]*Resource)
				}

				tcr := &Resource{
					Type: rs.Type,
				}

				tcr.Name = strings.TrimPrefix(crName, fmt.Sprintf("%s%s.", prefix, re.ResourceType))

				if cr.Change.Actions != nil {
					tcr.ChangeAction = Action(string(cr.Change.Actions[0]))

					if len(cr.Change.Actions) > 1 {
						tcr.ChangeAction = ActionReplace
					}
				}

				re.Children[crName] = tcr
			}

			if configured {

				var fname string
				ind := fmt.Sprintf("%s.%s", re.ResourceType, re.Name)

				if rs.Type == ResourceTypeResource && configs[parentConfig].Module.ManagedResources[ind] != nil {

					fname = filepath.Base(configs[parentConfig].Module.ManagedResources[ind].Pos.Filename)
					re.Line = &configs[parentConfig].Module.ManagedResources[ind].Pos.Line

					r.AddFileIfNotExists(parent, parentModule, fname)

					parent.Children[fname].Children[id] = re

				} else {

					r.AddFileIfNotExists(parent, parentModule, DefaultFileName)

					parent.Children[DefaultFileName].Children[id] = re
				}

			} else {

				parent.Children[id] = re
			}

		} else if rs.Type == ResourceTypeModule {
			re.Name = strings.Split(id, ".")[len(strings.Split(id, "."))-1]

			if configured && !childIndex.MatchString(id) && configs[parentConfig].Module.ModuleCalls[matchBrackets.ReplaceAllString(re.Name, "")] != nil {
				fname := filepath.Base(configs[parentConfig].Module.ModuleCalls[matchBrackets.ReplaceAllString(re.Name, "")].Pos.Filename)
				re.Line = &configs[parentConfig].Module.ModuleCalls[matchBrackets.ReplaceAllString(re.Name, "")].Pos.Line

				r.AddFileIfNotExists(parent, parentModule, fname)

				parent.Children[fname].Children[id] = re

			} else {
				parent.Children[id] = re
			}

			r.GenerateModuleMap(re, id)

		}
	}
}

func (r *rover) AddFileIfNotExists(module *Resource, parentModule string, fname string) {

	if _, ok := module.Children[fname]; !ok {

		module.Children[fname] = &Resource{
			Type:     ResourceTypeFile,
			Name:     fname,
			Source:   fmt.Sprintf("%s/%s", module.Source, fname),
			Children: map[string]*Resource{},
		}
	}
}

// Generates Map - Overview of files and their resources
// Groups different resource types together
// Defaults to config
func (r *rover) GenerateMap() error {
	log.Println("Generating resource map...")

	// Root module
	rootModule := &Resource{
		Type:     ResourceTypeModule,
		Name:     "",
		Source:   "unknown",
		Children: map[string]*Resource{},
	}

	mapObj := &Map{
		Path: "Rover Visualization",
		Root: rootModule.Children,
	}

	// If root module has local filesystem configuration stuff (line number/ file name info)
	rootConfig := r.RSO.Configs[""].Module

	if rootConfig != nil {
		rootModule.Source = rootConfig.Path
		mapObj.Path = rootConfig.Path
		mapObj.RequiredProviders = rootConfig.RequiredProviders
		mapObj.RequiredCore = rootConfig.RequiredCore
		r.GenerateModuleMap(rootModule, "")
	} else {
		r.AddFileIfNotExists(rootModule, "", DefaultFileName)
		r.GenerateModuleMap(rootModule.Children[DefaultFileName], "")
	}

	r.Map = mapObj

	return nil
}
