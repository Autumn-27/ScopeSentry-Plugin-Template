// Code generated by 'yaegi extract github.com/Autumn-27/ScopeSentry-Scan/internal/interfaces'. DO NOT EDIT.

package symbols

import (
	"github.com/Autumn-27/ScopeSentry-Scan/internal/interfaces"
	"reflect"
)

func init() {
	Symbols["github.com/Autumn-27/ScopeSentry-Scan/internal/interfaces/interfaces"] = map[string]reflect.Value{
		// type definitions
		"ModuleRunner": reflect.ValueOf((*interfaces.ModuleRunner)(nil)),
		"Plugin":       reflect.ValueOf((*interfaces.Plugin)(nil)),

		// interface wrapper definitions
		"_ModuleRunner": reflect.ValueOf((*_github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner)(nil)),
		"_Plugin":       reflect.ValueOf((*_github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin)(nil)),
	}
}

// _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner is an interface wrapper for ModuleRunner type
type _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner struct {
	IValue      interface{}
	WCloseInput func()
	WGetInput   func() chan interface{}
	WGetName    func() string
	WModuleRun  func() error
	WSetInput   func(a0 chan interface{})
}

func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner) CloseInput() {
	W.WCloseInput()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner) GetInput() chan interface{} {
	return W.WGetInput()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner) GetName() string {
	return W.WGetName()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner) ModuleRun() error {
	return W.WModuleRun()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_ModuleRunner) SetInput(a0 chan interface{}) {
	W.WSetInput(a0)
}

// _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin is an interface wrapper for Plugin type
type _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin struct {
	IValue        interface{}
	WCheck        func() error
	WClone        func() interfaces.Plugin
	WExecute      func(input interface{}) (interface{}, error)
	WGetCustom    func() interface{}
	WGetModule    func() string
	WGetName      func() string
	WGetParameter func() string
	WGetPluginId  func() string
	WGetTaskId    func() string
	WGetTaskName  func() string
	WInstall      func() error
	WLog          func(msg string, tp ...string)
	WSetCustom    func(cu interface{})
	WSetModule    func(name string)
	WSetName      func(name string)
	WSetParameter func(args string)
	WSetPluginId  func(id string)
	WSetResult    func(ch chan interface{})
	WSetTaskId    func(id string)
	WSetTaskName  func(name string)
	WUnInstall    func() error
}

func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) Check() error {
	return W.WCheck()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) Clone() interfaces.Plugin {
	return W.WClone()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) Execute(input interface{}) (interface{}, error) {
	return W.WExecute(input)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) GetCustom() interface{} {
	return W.WGetCustom()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) GetModule() string {
	return W.WGetModule()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) GetName() string {
	return W.WGetName()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) GetParameter() string {
	return W.WGetParameter()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) GetPluginId() string {
	return W.WGetPluginId()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) GetTaskId() string {
	return W.WGetTaskId()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) GetTaskName() string {
	return W.WGetTaskName()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) Install() error {
	return W.WInstall()
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) Log(msg string, tp ...string) {
	W.WLog(msg, tp...)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetCustom(cu interface{}) {
	W.WSetCustom(cu)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetModule(name string) {
	W.WSetModule(name)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetName(name string) {
	W.WSetName(name)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetParameter(args string) {
	W.WSetParameter(args)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetPluginId(id string) {
	W.WSetPluginId(id)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetResult(ch chan interface{}) {
	W.WSetResult(ch)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetTaskId(id string) {
	W.WSetTaskId(id)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) SetTaskName(name string) {
	W.WSetTaskName(name)
}
func (W _github_com_Autumn_27_ScopeSentry_Scan_internal_interfaces_Plugin) UnInstall() error {
	return W.WUnInstall()
}
