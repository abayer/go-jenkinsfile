package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Agent Determines the node/image in which the build will run from either named parameters or a bare none
type Agent struct {
	Argument  *RawArgument        `json:"argument,omitempty"`
	Arguments []*MapArgumentValue `json:"arguments,omitempty"`
	Type      string              `json:"type"`
}

// ArgumentList is a list of arguments or a single argument
type ArgumentList struct {
	Named      []*ArgumentValue
	Single     *RawArgument
	Positional []*RawArgument
}

// ArgumentValue The value for an argument
type ArgumentValue struct {
	Key   string       `json:"key,omitempty"`
	Value *RawArgument `json:"value,omitempty"`
}

// Axis One axis of a matrix
type Axis struct {
	Name   string         `json:"name"`
	Values []*RawArgument `json:"values"`
}

// Branch A block of steps, generally one of: the contents of a stage, the contents of a build condition block, or one branch of a parallel invocation
type Branch struct {
	Name  string     `json:"name"`
	Steps []*AnyStep `json:"steps"`
}

// AnyStep is either a step or a tree step
type AnyStep struct {
	Step *Step
	Tree *TreeStep
}

// BuildCondition A block of steps to be invoked depending on whether the given build condition is met
type BuildCondition struct {
	Branch    *Branch `json:"branch"`
	Condition string  `json:"condition"`
}

// EnvironmentEntry An entry in the environment
type EnvironmentEntry struct {
	Key   string            `json:"key,omitempty"`
	Value *EnvironmentValue `json:"value,omitempty"`
}

// EnvironmentValue is a value in the environment
type EnvironmentValue struct {
	Single   *RawArgument
	Function *InternalFunction
}

// ExcludeAxis One axis of a matrix
type ExcludeAxis struct {
	Inverse *bool          `json:"inverse,omitempty"`
	Name    *string        `json:"name"`
	Values  []*RawArgument `json:"values"`
}

// Input An input prompt for a stage
type Input struct {
	ID                 *RawArgument `json:"id,omitempty"`
	Message            *RawArgument `json:"message"`
	Ok                 *RawArgument `json:"ok,omitempty"`
	Parameters         *Parameters  `json:"parameters,omitempty"`
	Submitter          *RawArgument `json:"submitter,omitempty"`
	SubmitterParameter *RawArgument `json:"submitterParameter,omitempty"`
}

// InternalFunction An internal function call
type InternalFunction struct {
	Arguments []*RawArgument `json:"arguments,omitempty"`
	Name      string         `json:"name,omitempty"`
}

// KeyAndValueOrMethodCall A key/value pair that can either have a value or method call
type KeyAndValueOrMethodCall struct {
	Key   string             `json:"key,omitempty"`
	Value *ValueOrMethodCall `json:"value,omitempty"`
}

// ValueOrMethodCall is either a single value or a method call
type ValueOrMethodCall struct {
	Single *RawArgument
	Call   *MethodCall
}

// Libraries One or more shared library identifiers to load
type Libraries struct {
	Libraries []*RawArgument `json:"libraries,omitempty"`
}

// MapArgumentValue The value for a map argument
type MapArgumentValue struct {
	Key   string                     `json:"key,omitempty"`
	Value *MapArgumentValueRawOrList `json:"value,omitempty"`
}

// MapArgumentValueRawOrList is the raw argument or list of further arguments
type MapArgumentValueRawOrList struct {
	Raw  *RawArgument
	List []*MapArgumentValue
}

// Matrix Section containing a specification of a matrix - axes and stages
type Matrix struct {
	Agent       *Agent              `json:"agent,omitempty"`
	Axes        []*Axis             `json:"axes"`
	Environment []*EnvironmentEntry `json:"environment,omitempty"`
	Excludes    [][]*ExcludeAxis    `json:"excludes,omitempty"`
	Input       *Input              `json:"input,omitempty"`
	Options     *Options            `json:"options,omitempty"`
	Post        *Post               `json:"post,omitempty"`
	Stages      []*Stage            `json:"stages"`
	Tools       []*ArgumentValue    `json:"tools,omitempty"`
	When        *When               `json:"when,omitempty"`
}

// MethodArg is an argument to a method
type MethodArg struct {
	Single  *ValueOrMethodCall
	WithKey *KeyAndValueOrMethodCall
}

// MethodCall A method call with arguments, outside steps
type MethodCall struct {
	Arguments []*MethodArg `json:"arguments,omitempty"`
	Name      string       `json:"name,omitempty"`
}

// StepOrNestedWhenCondition is either a step or a nested when condition
type StepOrNestedWhenCondition struct {
	Step   *Step
	Nested *NestedWhenCondition
}

// NestedWhenCondition A when condition holding one or more other when conditions
type NestedWhenCondition struct {
	Children []*StepOrNestedWhenCondition `json:"children"`
	Name     string                       `json:"name"`
}

// Options One or more options (including job properties, wrappers, and options specific to Declarative Pipelines)
type Options struct {
	Options []*MethodCall `json:"options,omitempty"`
}

// Parameters One or more parameter definitions
type Parameters struct {
	Parameters []*MethodCall `json:"parameters,omitempty"`
}

// Pipeline defines the actual pipeline
type Pipeline struct {
	Agent       *Agent              `json:"agent"`
	Environment []*EnvironmentEntry `json:"environment,omitempty"`
	Libraries   *Libraries          `json:"libraries,omitempty"`
	Options     *Options            `json:"options,omitempty"`
	Parameters  *Parameters         `json:"parameters,omitempty"`
	Post        *Post               `json:"post,omitempty"`
	Stages      []*Stage            `json:"stages"`
	Tools       []*ArgumentValue    `json:"tools,omitempty"`
	Triggers    *Triggers           `json:"triggers,omitempty"`
}

// Post An array of build conditions with blocks of steps to run if those conditions are satisfied at the end of the build while still on the image/node the build ran on
type Post struct {
	Conditions []*BuildCondition `json:"conditions"`
}

// RawArgument The raw value of an argument, including whether it's a constant
type RawArgument struct {
	IsLiteral bool              `json:"isLiteral"`
	Value     *RawArgumentValue `json:"value"`
}

// RawArgumentValue is the value as one of a few possible types
type RawArgumentValue struct {
	AsFloat   *float64
	AsInteger *int64
	AsString  *string
	AsBool    *bool
}

// Root Schema for Kyoto AST JSON representation
type Root struct {
	Pipeline *Pipeline `json:"pipeline"`
}

// Stage A single Pipeline stage, with a name and either one or more branches or one or more nested stages
type Stage struct {
	Agent       *Agent              `json:"agent,omitempty"`
	Branches    []*Branch           `json:"branches,omitempty"`
	Environment []*EnvironmentEntry `json:"environment,omitempty"`
	FailFast    bool                `json:"failFast,omitempty"`
	Input       *Input              `json:"input,omitempty"`
	Matrix      *Matrix             `json:"matrix,omitempty"`
	Name        string              `json:"name"`
	Options     *Options            `json:"options,omitempty"`
	Parallel    []*Stage            `json:"parallel,omitempty"`
	Post        *Post               `json:"post,omitempty"`
	Stages      []*Stage            `json:"stages,omitempty"`
	Tools       []*ArgumentValue    `json:"tools,omitempty"`
	When        *When               `json:"when,omitempty"`
}

// Step A single step with parameters
type Step struct {
	Arguments *ArgumentList `json:"arguments"`
	Name      string        `json:"name"`
}

// TreeStep A block-scoped step with parameters containing 1 or more other steps
type TreeStep struct {
	Arguments *ArgumentList `json:"arguments"`
	Children  []*AnyStep    `json:"children"`
	Name      string        `json:"name"`
}

// Triggers One or more triggers
type Triggers struct {
	Triggers []*MethodCall `json:"triggers,omitempty"`
}

// When Conditions to evaluate whether the stage should run or not
type When struct {
	BeforeAgent   bool                         `json:"beforeAgent,omitempty"`
	BeforeInput   bool                         `json:"beforeInput,omitempty"`
	BeforeOptions bool                         `json:"beforeOptions,omitempty"`
	Conditions    []*StepOrNestedWhenCondition `json:"conditions"`
}

// MarshalJSON marshals the struct
func (strct *Agent) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "argument" field
	if comma {
		buf.WriteString(",")
	}
	if strct.Argument != nil {
		buf.WriteString("\"argument\": ")

		if tmp, err = json.Marshal(strct.Argument); err != nil {
			return nil, err
		}
		buf.Write(tmp)

		comma = true
	}
	// Marshal the "arguments" field
	if strct.Arguments != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"arguments\": ")
		if tmp, err = json.Marshal(strct.Arguments); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// "Type" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "type" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"type\": ")
	if tmp, err = json.Marshal(strct.Type); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Agent) UnmarshalJSON(b []byte) error {
	typeReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "argument":
			if err := json.Unmarshal([]byte(v), &strct.Argument); err != nil {
				return err
			}
		case "arguments":
			if err := json.Unmarshal([]byte(v), &strct.Arguments); err != nil {
				return err
			}
		case "type":
			if err := json.Unmarshal([]byte(v), &strct.Type); err != nil {
				return err
			}
			typeReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if type (a required property) was received
	if !typeReceived {
		return errors.New("\"type\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *ArgumentValue) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "key" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"key\": ")
	if tmp, err = json.Marshal(strct.Key); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "value" field
	if strct.Value != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"value\": ")
		if tmp, err = json.Marshal(strct.Value); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *ArgumentValue) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "key":
			if err := json.Unmarshal([]byte(v), &strct.Key); err != nil {
				return err
			}
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Axis) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Values" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "values" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"values\": ")
	if tmp, err = json.Marshal(strct.Values); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Axis) UnmarshalJSON(b []byte) error {
	nameReceived := false
	valuesReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		case "values":
			if err := json.Unmarshal([]byte(v), &strct.Values); err != nil {
				return err
			}
			valuesReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	// check if values (a required property) was received
	if !valuesReceived {
		return errors.New("\"values\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Branch) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Steps" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "steps" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"steps\": ")
	if tmp, err = json.Marshal(strct.Steps); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Branch) UnmarshalJSON(b []byte) error {
	nameReceived := false
	stepsReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		case "steps":
			if err := json.Unmarshal([]byte(v), &strct.Steps); err != nil {
				return err
			}
			stepsReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	// check if steps (a required property) was received
	if !stepsReceived {
		return errors.New("\"steps\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *BuildCondition) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Branch" field is required
	if strct.Branch == nil {
		return nil, errors.New("branch is a required field")
	}
	// Marshal the "branch" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"branch\": ")
	if tmp, err = json.Marshal(strct.Branch); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Condition" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "condition" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"condition\": ")
	if tmp, err = json.Marshal(strct.Condition); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *BuildCondition) UnmarshalJSON(b []byte) error {
	branchReceived := false
	conditionReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "branch":
			if err := json.Unmarshal([]byte(v), &strct.Branch); err != nil {
				return err
			}
			branchReceived = true
		case "condition":
			if err := json.Unmarshal([]byte(v), &strct.Condition); err != nil {
				return err
			}
			conditionReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if branch (a required property) was received
	if !branchReceived {
		return errors.New("\"branch\" is required but was not present")
	}
	// check if condition (a required property) was received
	if !conditionReceived {
		return errors.New("\"condition\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *EnvironmentEntry) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "key" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"key\": ")
	if tmp, err = json.Marshal(strct.Key); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "value" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"value\": ")
	if tmp, err = json.Marshal(strct.Value); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *EnvironmentEntry) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "key":
			if err := json.Unmarshal([]byte(v), &strct.Key); err != nil {
				return err
			}
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *ExcludeAxis) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "inverse" field
	if strct.Inverse != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"inverse\": ")
		if tmp, err = json.Marshal(strct.Inverse); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Values" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "values" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"values\": ")
	if tmp, err = json.Marshal(strct.Values); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *ExcludeAxis) UnmarshalJSON(b []byte) error {
	nameReceived := false
	valuesReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "inverse":
			if err := json.Unmarshal([]byte(v), &strct.Inverse); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		case "values":
			if err := json.Unmarshal([]byte(v), &strct.Values); err != nil {
				return err
			}
			valuesReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	// check if values (a required property) was received
	if !valuesReceived {
		return errors.New("\"values\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Input) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "id" field
	if comma {
		buf.WriteString(",")
	}
	if strct.ID != nil {
		buf.WriteString("\"id\": ")
		if tmp, err = json.Marshal(strct.ID); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// "Message" field is required
	if strct.Message == nil {
		return nil, errors.New("message is a required field")
	}
	// Marshal the "message" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"message\": ")
	if tmp, err = json.Marshal(strct.Message); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "ok" field
	if strct.Ok != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"ok\": ")
		if tmp, err = json.Marshal(strct.Ok); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "parameters" field
	if strct.Parameters != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"parameters\": ")
		if tmp, err = json.Marshal(strct.Parameters); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "submitter" field
	if strct.Submitter != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"submitter\": ")
		if tmp, err = json.Marshal(strct.Submitter); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "submitterParameter" field
	if strct.SubmitterParameter != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"submitterParameter\": ")
		if tmp, err = json.Marshal(strct.SubmitterParameter); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Input) UnmarshalJSON(b []byte) error {
	messageReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "id":
			if err := json.Unmarshal([]byte(v), &strct.ID); err != nil {
				return err
			}
		case "message":
			if err := json.Unmarshal([]byte(v), &strct.Message); err != nil {
				return err
			}
			messageReceived = true
		case "ok":
			if err := json.Unmarshal([]byte(v), &strct.Ok); err != nil {
				return err
			}
		case "parameters":
			if err := json.Unmarshal([]byte(v), &strct.Parameters); err != nil {
				return err
			}
		case "submitter":
			if err := json.Unmarshal([]byte(v), &strct.Submitter); err != nil {
				return err
			}
		case "submitterParameter":
			if err := json.Unmarshal([]byte(v), &strct.SubmitterParameter); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if message (a required property) was received
	if !messageReceived {
		return errors.New("\"message\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *InternalFunction) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "arguments" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"arguments\": ")
	if tmp, err = json.Marshal(strct.Arguments); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *InternalFunction) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "arguments":
			if err := json.Unmarshal([]byte(v), &strct.Arguments); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *KeyAndValueOrMethodCall) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "key" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"key\": ")
	if tmp, err = json.Marshal(strct.Key); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "value" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"value\": ")
	if tmp, err = json.Marshal(strct.Value); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *KeyAndValueOrMethodCall) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "key":
			if err := json.Unmarshal([]byte(v), &strct.Key); err != nil {
				return err
			}
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Libraries) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "libraries" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"libraries\": ")
	if tmp, err = json.Marshal(strct.Libraries); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Libraries) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "libraries":
			if err := json.Unmarshal([]byte(v), &strct.Libraries); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *MapArgumentValue) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "key" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"key\": ")
	if tmp, err = json.Marshal(strct.Key); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "value" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"value\": ")
	if tmp, err = json.Marshal(strct.Value); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *MapArgumentValue) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "key":
			if err := json.Unmarshal([]byte(v), &strct.Key); err != nil {
				return err
			}
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *MapArgumentValueRawOrList) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	if strct.Raw != nil {
		return strct.Raw.MarshalJSON()
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	if tmp, err = json.Marshal(strct.List); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *MapArgumentValueRawOrList) UnmarshalJSON(b []byte) error {
	var err error

	if err = json.Unmarshal(b, &strct.List); err == nil {
		return nil
	}
	err = json.Unmarshal(b, &strct.Raw)

	return err
}

// MarshalJSON marshals the struct
func (strct *ArgumentList) MarshalJSON() ([]byte, error) {
	if strct.Single != nil {
		return strct.Single.MarshalJSON()
	}

	var err error
	var tmp []byte
	if strct.Named != nil {
		tmp, err = json.Marshal(strct.Named)
	} else if strct.Positional != nil {
		tmp, err = json.Marshal(strct.Positional)
	}
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *ArgumentList) UnmarshalJSON(b []byte) error {
	var err error

	if err = json.Unmarshal(b, &strct.Named); err == nil {
		return nil
	}
	if err = json.Unmarshal(b, &strct.Positional); err == nil {
		return nil
	}
	if err = json.Unmarshal(b, &strct.Single); err == nil {
		return nil
	}

	return err
}

// MarshalJSON marshals the struct
func (strct *AnyStep) MarshalJSON() ([]byte, error) {
	if strct.Step != nil {
		return strct.Step.MarshalJSON()
	}
	if strct.Tree != nil {
		return strct.Tree.MarshalJSON()
	}
	return nil, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *AnyStep) UnmarshalJSON(b []byte) error {
	var err error

	if err = json.Unmarshal(b, &strct.Tree); err == nil {
		strct.Step = nil
		return nil
	}
	if err = json.Unmarshal(b, &strct.Step); err == nil {
		strct.Tree = nil
		return nil
	}

	return err
}

// MarshalJSON marshals the struct
func (strct *EnvironmentValue) MarshalJSON() ([]byte, error) {
	if strct.Function != nil {
		return strct.Function.MarshalJSON()
	}
	if strct.Single != nil {
		return strct.Single.MarshalJSON()
	}
	return nil, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *EnvironmentValue) UnmarshalJSON(b []byte) error {
	var err error

	if err = json.Unmarshal(b, &strct.Function); err == nil {
		strct.Single = nil
		return nil
	}
	if err = json.Unmarshal(b, &strct.Single); err == nil {
		strct.Function = nil
		return nil
	}

	return err
}

// MarshalJSON marshals the struct
func (strct *ValueOrMethodCall) MarshalJSON() ([]byte, error) {
	if strct.Call != nil {
		return strct.Call.MarshalJSON()
	}
	if strct.Single != nil {
		return strct.Single.MarshalJSON()
	}
	return nil, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *ValueOrMethodCall) UnmarshalJSON(b []byte) error {
	var err error

	if err = json.Unmarshal(b, &strct.Call); err == nil {
		return nil
	}
	if err = json.Unmarshal(b, &strct.Single); err == nil {
		return nil
	}

	return err
}

// MarshalJSON marshals the struct
func (strct *MethodArg) MarshalJSON() ([]byte, error) {
	if strct.Single != nil {
		return strct.Single.MarshalJSON()
	}
	if strct.WithKey != nil {
		return strct.WithKey.MarshalJSON()
	}
	return nil, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *MethodArg) UnmarshalJSON(b []byte) error {
	var err error

	if err = json.Unmarshal(b, &strct.WithKey); err == nil {
		return nil
	}
	if err = json.Unmarshal(b, &strct.Single); err == nil {
		return nil
	}

	return err
}

// MarshalJSON marshals the struct
func (strct *StepOrNestedWhenCondition) MarshalJSON() ([]byte, error) {
	if strct.Step != nil {
		return strct.Step.MarshalJSON()
	}
	if strct.Nested != nil {
		return strct.Nested.MarshalJSON()
	}
	return nil, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *StepOrNestedWhenCondition) UnmarshalJSON(b []byte) error {
	var err error

	if err = json.Unmarshal(b, &strct.Nested); err == nil {
		return nil
	}
	if err = json.Unmarshal(b, &strct.Step); err == nil {
		return nil
	}

	return err
}

// MarshalJSON marshals the struct
func (strct *Matrix) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "agent" field
	if strct.Agent != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"agent\": ")
		if tmp, err = json.Marshal(strct.Agent); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// "Axes" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "axes" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"axes\": ")
	if tmp, err = json.Marshal(strct.Axes); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "environment" field
	if strct.Environment != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"environment\": ")
		if tmp, err = json.Marshal(strct.Environment); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "excludes" field
	if strct.Excludes != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"excludes\": ")
		if tmp, err = json.Marshal(strct.Excludes); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "input" field
	if strct.Input != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"input\": ")
		if tmp, err = json.Marshal(strct.Input); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "options" field
	if strct.Options != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"options\": ")
		if tmp, err = json.Marshal(strct.Options); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "post" field
	if strct.Post != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"post\": ")
		if tmp, err = json.Marshal(strct.Post); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// "Stages" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "stages" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"stages\": ")
	if tmp, err = json.Marshal(strct.Stages); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "tools" field
	if strct.Tools != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"tools\": ")
		if tmp, err = json.Marshal(strct.Tools); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "when" field
	if strct.When != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"when\": ")
		if tmp, err = json.Marshal(strct.When); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Matrix) UnmarshalJSON(b []byte) error {
	axesReceived := false
	stagesReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "agent":
			if err := json.Unmarshal([]byte(v), &strct.Agent); err != nil {
				return err
			}
		case "axes":
			if err := json.Unmarshal([]byte(v), &strct.Axes); err != nil {
				return err
			}
			axesReceived = true
		case "environment":
			if err := json.Unmarshal([]byte(v), &strct.Environment); err != nil {
				return err
			}
		case "excludes":
			if err := json.Unmarshal([]byte(v), &strct.Excludes); err != nil {
				return err
			}
		case "input":
			if err := json.Unmarshal([]byte(v), &strct.Input); err != nil {
				return err
			}
		case "options":
			if err := json.Unmarshal([]byte(v), &strct.Options); err != nil {
				return err
			}
		case "post":
			if err := json.Unmarshal([]byte(v), &strct.Post); err != nil {
				return err
			}
		case "stages":
			if err := json.Unmarshal([]byte(v), &strct.Stages); err != nil {
				return err
			}
			stagesReceived = true
		case "tools":
			if err := json.Unmarshal([]byte(v), &strct.Tools); err != nil {
				return err
			}
		case "when":
			if err := json.Unmarshal([]byte(v), &strct.When); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if axes (a required property) was received
	if !axesReceived {
		return errors.New("\"axes\" is required but was not present")
	}
	// check if stages (a required property) was received
	if !stagesReceived {
		return errors.New("\"stages\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *MethodCall) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "arguments" field
	if len(strct.Arguments) > 0 {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"arguments\": ")
		if tmp, err = json.Marshal(strct.Arguments); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *MethodCall) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "arguments":
			if err := json.Unmarshal([]byte(v), &strct.Arguments); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *NestedWhenCondition) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Children" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "children" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"children\": ")
	if tmp, err = json.Marshal(strct.Children); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *NestedWhenCondition) UnmarshalJSON(b []byte) error {
	childrenReceived := false
	nameReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "children":
			if err := json.Unmarshal([]byte(v), &strct.Children); err != nil {
				return err
			}
			childrenReceived = true
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if children (a required property) was received
	if !childrenReceived {
		return errors.New("\"children\" is required but was not present")
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Options) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "options" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"options\": ")
	if tmp, err = json.Marshal(strct.Options); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Options) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "options":
			if err := json.Unmarshal([]byte(v), &strct.Options); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Parameters) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "parameters" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"parameters\": ")
	if tmp, err = json.Marshal(strct.Parameters); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Parameters) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "parameters":
			if err := json.Unmarshal([]byte(v), &strct.Parameters); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Pipeline) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Agent" field is required
	if strct.Agent == nil {
		return nil, errors.New("agent is a required field")
	}
	// Marshal the "agent" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"agent\": ")
	if tmp, err = json.Marshal(strct.Agent); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "environment" field
	if strct.Environment != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"environment\": ")
		if tmp, err = json.Marshal(strct.Environment); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "libraries" field
	if strct.Libraries != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"libraries\": ")
		if tmp, err = json.Marshal(strct.Libraries); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "options" field
	if strct.Options != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"options\": ")
		if tmp, err = json.Marshal(strct.Options); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "parameters" field
	if strct.Parameters != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"parameters\": ")
		if tmp, err = json.Marshal(strct.Parameters); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "post" field
	if strct.Post != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"post\": ")
		if tmp, err = json.Marshal(strct.Post); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// "Stages" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "stages" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"stages\": ")
	if tmp, err = json.Marshal(strct.Stages); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "tools" field
	if strct.Tools != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"tools\": ")
		if tmp, err = json.Marshal(strct.Tools); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "triggers" field
	if strct.Triggers != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"triggers\": ")
		if tmp, err = json.Marshal(strct.Triggers); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Pipeline) UnmarshalJSON(b []byte) error {
	agentReceived := false
	stagesReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "agent":
			if err := json.Unmarshal([]byte(v), &strct.Agent); err != nil {
				return err
			}
			agentReceived = true
		case "environment":
			if err := json.Unmarshal([]byte(v), &strct.Environment); err != nil {
				return err
			}
		case "libraries":
			if err := json.Unmarshal([]byte(v), &strct.Libraries); err != nil {
				return err
			}
		case "options":
			if err := json.Unmarshal([]byte(v), &strct.Options); err != nil {
				return err
			}
		case "parameters":
			if err := json.Unmarshal([]byte(v), &strct.Parameters); err != nil {
				return err
			}
		case "post":
			if err := json.Unmarshal([]byte(v), &strct.Post); err != nil {
				return err
			}
		case "stages":
			if err := json.Unmarshal([]byte(v), &strct.Stages); err != nil {
				return err
			}
			stagesReceived = true
		case "tools":
			if err := json.Unmarshal([]byte(v), &strct.Tools); err != nil {
				return err
			}
		case "triggers":
			if err := json.Unmarshal([]byte(v), &strct.Triggers); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if agent (a required property) was received
	if !agentReceived {
		return errors.New("\"agent\" is required but was not present")
	}
	// check if stages (a required property) was received
	if !stagesReceived {
		return errors.New("\"stages\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Post) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Conditions" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "conditions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"conditions\": ")
	if tmp, err = json.Marshal(strct.Conditions); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Post) UnmarshalJSON(b []byte) error {
	conditionsReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "conditions":
			if err := json.Unmarshal([]byte(v), &strct.Conditions); err != nil {
				return err
			}
			conditionsReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if conditions (a required property) was received
	if !conditionsReceived {
		return errors.New("\"conditions\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *RawArgument) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "IsLiteral" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "isLiteral" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"isLiteral\": ")
	if tmp, err = json.Marshal(strct.IsLiteral); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Value" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "value" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"value\": ")
	if tmp, err = json.Marshal(strct.Value); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *RawArgument) UnmarshalJSON(b []byte) error {
	isLiteralReceived := false
	valueReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "isLiteral":
			if err := json.Unmarshal([]byte(v), &strct.IsLiteral); err != nil {
				return err
			}
			isLiteralReceived = true
		case "value":
			if err := json.Unmarshal([]byte(v), &strct.Value); err != nil {
				return err
			}
			valueReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if isLiteral (a required property) was received
	if !isLiteralReceived {
		return errors.New("\"isLiteral\" is required but was not present")
	}
	// check if value (a required property) was received
	if !valueReceived {
		return errors.New("\"value\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *RawArgumentValue) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	if strct.AsBool != nil {
		buf.WriteString(fmt.Sprintf("%t", *strct.AsBool))
	} else if strct.AsFloat != nil {
		buf.WriteString(fmt.Sprintf("%f", *strct.AsFloat))
	} else if strct.AsInteger != nil {
		buf.WriteString(fmt.Sprintf("%d", *strct.AsInteger))
	} else if strct.AsString != nil {
		buf.WriteString(strconv.Quote(*strct.AsString))
	} else {
		buf.WriteString("\"\"")
	}

	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *RawArgumentValue) UnmarshalJSON(b []byte) error {
	var err error
	if err = json.Unmarshal(b, &strct.AsBool); err == nil {
		return nil
	}
	strct.AsBool = nil
	if err = json.Unmarshal(b, &strct.AsFloat); err == nil {
		return nil
	}
	strct.AsFloat = nil
	if err = json.Unmarshal(b, &strct.AsInteger); err == nil {
		return nil
	}
	strct.AsInteger = nil
	if err = json.Unmarshal(b, &strct.AsString); err == nil {
		return nil
	}
	return err
}

// MarshalJSON marshals the struct
func (strct *Root) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Pipeline" field is required
	if strct.Pipeline == nil {
		return nil, errors.New("pipeline is a required field")
	}
	// Marshal the "pipeline" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"pipeline\": ")
	if tmp, err = json.Marshal(strct.Pipeline); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Root) UnmarshalJSON(b []byte) error {
	pipelineReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "pipeline":
			if err := json.Unmarshal([]byte(v), &strct.Pipeline); err != nil {
				return err
			}
			pipelineReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if pipeline (a required property) was received
	if !pipelineReceived {
		return errors.New("\"pipeline\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Stage) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "agent" field
	if strct.Agent != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"agent\": ")
		if tmp, err = json.Marshal(strct.Agent); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "branches" field
	if strct.Branches != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"branches\": ")
		if tmp, err = json.Marshal(strct.Branches); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "environment" field
	if strct.Environment != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"environment\": ")
		if tmp, err = json.Marshal(strct.Environment); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "failFast" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"failFast\": ")
	if tmp, err = json.Marshal(strct.FailFast); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "input" field
	if strct.Input != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"input\": ")
		if tmp, err = json.Marshal(strct.Input); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "matrix" field
	if strct.Matrix != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"matrix\": ")
		if tmp, err = json.Marshal(strct.Matrix); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "options" field
	if strct.Options != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"options\": ")
		if tmp, err = json.Marshal(strct.Options); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "parallel" field
	if strct.Parallel != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"parallel\": ")
		if tmp, err = json.Marshal(strct.Parallel); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "post" field
	if strct.Post != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"post\": ")
		if tmp, err = json.Marshal(strct.Post); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "stages" field
	if strct.Stages != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"stages\": ")
		if tmp, err = json.Marshal(strct.Stages); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "tools" field
	if strct.Tools != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"tools\": ")
		if tmp, err = json.Marshal(strct.Tools); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}
	// Marshal the "when" field
	if strct.When != nil {
		if comma {
			buf.WriteString(",")
		}
		buf.WriteString("\"when\": ")
		if tmp, err = json.Marshal(strct.When); err != nil {
			return nil, err
		}
		buf.Write(tmp)
		comma = true
	}

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Stage) UnmarshalJSON(b []byte) error {
	nameReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "agent":
			if err := json.Unmarshal([]byte(v), &strct.Agent); err != nil {
				return err
			}
		case "branches":
			if err := json.Unmarshal([]byte(v), &strct.Branches); err != nil {
				return err
			}
		case "environment":
			if err := json.Unmarshal([]byte(v), &strct.Environment); err != nil {
				return err
			}
		case "failFast":
			if err := json.Unmarshal([]byte(v), &strct.FailFast); err != nil {
				return err
			}
		case "input":
			if err := json.Unmarshal([]byte(v), &strct.Input); err != nil {
				return err
			}
		case "matrix":
			if err := json.Unmarshal([]byte(v), &strct.Matrix); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		case "options":
			if err := json.Unmarshal([]byte(v), &strct.Options); err != nil {
				return err
			}
		case "parallel":
			if err := json.Unmarshal([]byte(v), &strct.Parallel); err != nil {
				return err
			}
		case "post":
			if err := json.Unmarshal([]byte(v), &strct.Post); err != nil {
				return err
			}
		case "stages":
			if err := json.Unmarshal([]byte(v), &strct.Stages); err != nil {
				return err
			}
		case "tools":
			if err := json.Unmarshal([]byte(v), &strct.Tools); err != nil {
				return err
			}
		case "when":
			if err := json.Unmarshal([]byte(v), &strct.When); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Step) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Arguments" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "arguments" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"arguments\": ")
	if tmp, err = json.Marshal(strct.Arguments); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Step) UnmarshalJSON(b []byte) error {
	argumentsReceived := false
	nameReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "arguments":
			if err := json.Unmarshal([]byte(v), &strct.Arguments); err != nil {
				return err
			}
			argumentsReceived = true
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if arguments (a required property) was received
	if !argumentsReceived {
		return errors.New("\"arguments\" is required but was not present")
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *TreeStep) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Arguments" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "arguments" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"arguments\": ")
	if tmp, err = json.Marshal(strct.Arguments); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Children" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "children" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"children\": ")
	if tmp, err = json.Marshal(strct.Children); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err = json.Marshal(strct.Name); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *TreeStep) UnmarshalJSON(b []byte) error {
	argumentsReceived := false
	childrenReceived := false
	nameReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "arguments":
			if err := json.Unmarshal([]byte(v), &strct.Arguments); err != nil {
				return err
			}
			argumentsReceived = true
		case "children":
			if err := json.Unmarshal([]byte(v), &strct.Children); err != nil {
				return err
			}
			childrenReceived = true
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if arguments (a required property) was received
	if !argumentsReceived {
		return errors.New("\"arguments\" is required but was not present")
	}
	// check if children (a required property) was received
	if !childrenReceived {
		return errors.New("\"children\" is required but was not present")
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *Triggers) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "triggers" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"triggers\": ")
	if tmp, err = json.Marshal(strct.Triggers); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *Triggers) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "triggers":
			if err := json.Unmarshal([]byte(v), &strct.Triggers); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

// MarshalJSON marshals the struct
func (strct *When) MarshalJSON() ([]byte, error) {
	var tmp []byte
	var err error
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "beforeAgent" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"beforeAgent\": ")
	if tmp, err = json.Marshal(strct.BeforeAgent); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "beforeInput" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"beforeInput\": ")
	if tmp, err = json.Marshal(strct.BeforeInput); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// Marshal the "beforeOptions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"beforeOptions\": ")
	if tmp, err = json.Marshal(strct.BeforeOptions); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true
	// "Conditions" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "conditions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"conditions\": ")
	if tmp, err = json.Marshal(strct.Conditions); err != nil {
		return nil, err
	}
	buf.Write(tmp)
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

// UnmarshalJSON unmarshals the struct
func (strct *When) UnmarshalJSON(b []byte) error {
	conditionsReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "beforeAgent":
			if err := json.Unmarshal([]byte(v), &strct.BeforeAgent); err != nil {
				return err
			}
		case "beforeInput":
			if err := json.Unmarshal([]byte(v), &strct.BeforeInput); err != nil {
				return err
			}
		case "beforeOptions":
			if err := json.Unmarshal([]byte(v), &strct.BeforeOptions); err != nil {
				return err
			}
		case "conditions":
			if err := json.Unmarshal([]byte(v), &strct.Conditions); err != nil {
				return err
			}
			conditionsReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if conditions (a required property) was received
	if !conditionsReceived {
		return errors.New("\"conditions\" is required but was not present")
	}
	return nil
}
