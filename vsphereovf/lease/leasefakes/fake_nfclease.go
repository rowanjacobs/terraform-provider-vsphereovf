// Code generated by counterfeiter. DO NOT EDIT.
package leasefakes

import (
	"context"
	"io"
	"sync"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/lease"
	"github.com/vmware/govmomi/nfc"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
)

type FakeNFCLease struct {
	UploadStub        func(context.Context, nfc.FileItem, io.Reader, soap.Upload) error
	uploadMutex       sync.RWMutex
	uploadArgsForCall []struct {
		arg1 context.Context
		arg2 nfc.FileItem
		arg3 io.Reader
		arg4 soap.Upload
	}
	uploadReturns struct {
		result1 error
	}
	uploadReturnsOnCall map[int]struct {
		result1 error
	}
	WaitStub        func(context.Context, []types.OvfFileItem) (*nfc.LeaseInfo, error)
	waitMutex       sync.RWMutex
	waitArgsForCall []struct {
		arg1 context.Context
		arg2 []types.OvfFileItem
	}
	waitReturns struct {
		result1 *nfc.LeaseInfo
		result2 error
	}
	waitReturnsOnCall map[int]struct {
		result1 *nfc.LeaseInfo
		result2 error
	}
	StartUpdaterStub        func(context.Context, *nfc.LeaseInfo) *nfc.LeaseUpdater
	startUpdaterMutex       sync.RWMutex
	startUpdaterArgsForCall []struct {
		arg1 context.Context
		arg2 *nfc.LeaseInfo
	}
	startUpdaterReturns struct {
		result1 *nfc.LeaseUpdater
	}
	startUpdaterReturnsOnCall map[int]struct {
		result1 *nfc.LeaseUpdater
	}
	CompleteStub        func(context.Context) error
	completeMutex       sync.RWMutex
	completeArgsForCall []struct {
		arg1 context.Context
	}
	completeReturns struct {
		result1 error
	}
	completeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeNFCLease) Upload(arg1 context.Context, arg2 nfc.FileItem, arg3 io.Reader, arg4 soap.Upload) error {
	fake.uploadMutex.Lock()
	ret, specificReturn := fake.uploadReturnsOnCall[len(fake.uploadArgsForCall)]
	fake.uploadArgsForCall = append(fake.uploadArgsForCall, struct {
		arg1 context.Context
		arg2 nfc.FileItem
		arg3 io.Reader
		arg4 soap.Upload
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("Upload", []interface{}{arg1, arg2, arg3, arg4})
	fake.uploadMutex.Unlock()
	if fake.UploadStub != nil {
		return fake.UploadStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.uploadReturns.result1
}

func (fake *FakeNFCLease) UploadCallCount() int {
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	return len(fake.uploadArgsForCall)
}

func (fake *FakeNFCLease) UploadArgsForCall(i int) (context.Context, nfc.FileItem, io.Reader, soap.Upload) {
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	return fake.uploadArgsForCall[i].arg1, fake.uploadArgsForCall[i].arg2, fake.uploadArgsForCall[i].arg3, fake.uploadArgsForCall[i].arg4
}

func (fake *FakeNFCLease) UploadReturns(result1 error) {
	fake.UploadStub = nil
	fake.uploadReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNFCLease) UploadReturnsOnCall(i int, result1 error) {
	fake.UploadStub = nil
	if fake.uploadReturnsOnCall == nil {
		fake.uploadReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.uploadReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeNFCLease) Wait(arg1 context.Context, arg2 []types.OvfFileItem) (*nfc.LeaseInfo, error) {
	var arg2Copy []types.OvfFileItem
	if arg2 != nil {
		arg2Copy = make([]types.OvfFileItem, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.waitMutex.Lock()
	ret, specificReturn := fake.waitReturnsOnCall[len(fake.waitArgsForCall)]
	fake.waitArgsForCall = append(fake.waitArgsForCall, struct {
		arg1 context.Context
		arg2 []types.OvfFileItem
	}{arg1, arg2Copy})
	fake.recordInvocation("Wait", []interface{}{arg1, arg2Copy})
	fake.waitMutex.Unlock()
	if fake.WaitStub != nil {
		return fake.WaitStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.waitReturns.result1, fake.waitReturns.result2
}

func (fake *FakeNFCLease) WaitCallCount() int {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return len(fake.waitArgsForCall)
}

func (fake *FakeNFCLease) WaitArgsForCall(i int) (context.Context, []types.OvfFileItem) {
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	return fake.waitArgsForCall[i].arg1, fake.waitArgsForCall[i].arg2
}

func (fake *FakeNFCLease) WaitReturns(result1 *nfc.LeaseInfo, result2 error) {
	fake.WaitStub = nil
	fake.waitReturns = struct {
		result1 *nfc.LeaseInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeNFCLease) WaitReturnsOnCall(i int, result1 *nfc.LeaseInfo, result2 error) {
	fake.WaitStub = nil
	if fake.waitReturnsOnCall == nil {
		fake.waitReturnsOnCall = make(map[int]struct {
			result1 *nfc.LeaseInfo
			result2 error
		})
	}
	fake.waitReturnsOnCall[i] = struct {
		result1 *nfc.LeaseInfo
		result2 error
	}{result1, result2}
}

func (fake *FakeNFCLease) StartUpdater(arg1 context.Context, arg2 *nfc.LeaseInfo) *nfc.LeaseUpdater {
	fake.startUpdaterMutex.Lock()
	ret, specificReturn := fake.startUpdaterReturnsOnCall[len(fake.startUpdaterArgsForCall)]
	fake.startUpdaterArgsForCall = append(fake.startUpdaterArgsForCall, struct {
		arg1 context.Context
		arg2 *nfc.LeaseInfo
	}{arg1, arg2})
	fake.recordInvocation("StartUpdater", []interface{}{arg1, arg2})
	fake.startUpdaterMutex.Unlock()
	if fake.StartUpdaterStub != nil {
		return fake.StartUpdaterStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.startUpdaterReturns.result1
}

func (fake *FakeNFCLease) StartUpdaterCallCount() int {
	fake.startUpdaterMutex.RLock()
	defer fake.startUpdaterMutex.RUnlock()
	return len(fake.startUpdaterArgsForCall)
}

func (fake *FakeNFCLease) StartUpdaterArgsForCall(i int) (context.Context, *nfc.LeaseInfo) {
	fake.startUpdaterMutex.RLock()
	defer fake.startUpdaterMutex.RUnlock()
	return fake.startUpdaterArgsForCall[i].arg1, fake.startUpdaterArgsForCall[i].arg2
}

func (fake *FakeNFCLease) StartUpdaterReturns(result1 *nfc.LeaseUpdater) {
	fake.StartUpdaterStub = nil
	fake.startUpdaterReturns = struct {
		result1 *nfc.LeaseUpdater
	}{result1}
}

func (fake *FakeNFCLease) StartUpdaterReturnsOnCall(i int, result1 *nfc.LeaseUpdater) {
	fake.StartUpdaterStub = nil
	if fake.startUpdaterReturnsOnCall == nil {
		fake.startUpdaterReturnsOnCall = make(map[int]struct {
			result1 *nfc.LeaseUpdater
		})
	}
	fake.startUpdaterReturnsOnCall[i] = struct {
		result1 *nfc.LeaseUpdater
	}{result1}
}

func (fake *FakeNFCLease) Complete(arg1 context.Context) error {
	fake.completeMutex.Lock()
	ret, specificReturn := fake.completeReturnsOnCall[len(fake.completeArgsForCall)]
	fake.completeArgsForCall = append(fake.completeArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	fake.recordInvocation("Complete", []interface{}{arg1})
	fake.completeMutex.Unlock()
	if fake.CompleteStub != nil {
		return fake.CompleteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.completeReturns.result1
}

func (fake *FakeNFCLease) CompleteCallCount() int {
	fake.completeMutex.RLock()
	defer fake.completeMutex.RUnlock()
	return len(fake.completeArgsForCall)
}

func (fake *FakeNFCLease) CompleteArgsForCall(i int) context.Context {
	fake.completeMutex.RLock()
	defer fake.completeMutex.RUnlock()
	return fake.completeArgsForCall[i].arg1
}

func (fake *FakeNFCLease) CompleteReturns(result1 error) {
	fake.CompleteStub = nil
	fake.completeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeNFCLease) CompleteReturnsOnCall(i int, result1 error) {
	fake.CompleteStub = nil
	if fake.completeReturnsOnCall == nil {
		fake.completeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.completeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeNFCLease) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	fake.waitMutex.RLock()
	defer fake.waitMutex.RUnlock()
	fake.startUpdaterMutex.RLock()
	defer fake.startUpdaterMutex.RUnlock()
	fake.completeMutex.RLock()
	defer fake.completeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeNFCLease) recordInvocation(key string, args []interface{}) {
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

var _ lease.NFCLease = new(FakeNFCLease)
