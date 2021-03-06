// Code generated by counterfeiter. DO NOT EDIT.
package importerfakes

import (
	"context"
	"sync"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/importer"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type FakeOVFManager struct {
	CreateImportSpecStub        func(context.Context, string, mo.Reference, mo.Reference, types.OvfCreateImportSpecParams) (*types.OvfCreateImportSpecResult, error)
	createImportSpecMutex       sync.RWMutex
	createImportSpecArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 mo.Reference
		arg4 mo.Reference
		arg5 types.OvfCreateImportSpecParams
	}
	createImportSpecReturns struct {
		result1 *types.OvfCreateImportSpecResult
		result2 error
	}
	createImportSpecReturnsOnCall map[int]struct {
		result1 *types.OvfCreateImportSpecResult
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOVFManager) CreateImportSpec(arg1 context.Context, arg2 string, arg3 mo.Reference, arg4 mo.Reference, arg5 types.OvfCreateImportSpecParams) (*types.OvfCreateImportSpecResult, error) {
	fake.createImportSpecMutex.Lock()
	ret, specificReturn := fake.createImportSpecReturnsOnCall[len(fake.createImportSpecArgsForCall)]
	fake.createImportSpecArgsForCall = append(fake.createImportSpecArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 mo.Reference
		arg4 mo.Reference
		arg5 types.OvfCreateImportSpecParams
	}{arg1, arg2, arg3, arg4, arg5})
	fake.recordInvocation("CreateImportSpec", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.createImportSpecMutex.Unlock()
	if fake.CreateImportSpecStub != nil {
		return fake.CreateImportSpecStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createImportSpecReturns.result1, fake.createImportSpecReturns.result2
}

func (fake *FakeOVFManager) CreateImportSpecCallCount() int {
	fake.createImportSpecMutex.RLock()
	defer fake.createImportSpecMutex.RUnlock()
	return len(fake.createImportSpecArgsForCall)
}

func (fake *FakeOVFManager) CreateImportSpecArgsForCall(i int) (context.Context, string, mo.Reference, mo.Reference, types.OvfCreateImportSpecParams) {
	fake.createImportSpecMutex.RLock()
	defer fake.createImportSpecMutex.RUnlock()
	return fake.createImportSpecArgsForCall[i].arg1, fake.createImportSpecArgsForCall[i].arg2, fake.createImportSpecArgsForCall[i].arg3, fake.createImportSpecArgsForCall[i].arg4, fake.createImportSpecArgsForCall[i].arg5
}

func (fake *FakeOVFManager) CreateImportSpecReturns(result1 *types.OvfCreateImportSpecResult, result2 error) {
	fake.CreateImportSpecStub = nil
	fake.createImportSpecReturns = struct {
		result1 *types.OvfCreateImportSpecResult
		result2 error
	}{result1, result2}
}

func (fake *FakeOVFManager) CreateImportSpecReturnsOnCall(i int, result1 *types.OvfCreateImportSpecResult, result2 error) {
	fake.CreateImportSpecStub = nil
	if fake.createImportSpecReturnsOnCall == nil {
		fake.createImportSpecReturnsOnCall = make(map[int]struct {
			result1 *types.OvfCreateImportSpecResult
			result2 error
		})
	}
	fake.createImportSpecReturnsOnCall[i] = struct {
		result1 *types.OvfCreateImportSpecResult
		result2 error
	}{result1, result2}
}

func (fake *FakeOVFManager) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createImportSpecMutex.RLock()
	defer fake.createImportSpecMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOVFManager) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ importer.OVFManager = new(FakeOVFManager)
