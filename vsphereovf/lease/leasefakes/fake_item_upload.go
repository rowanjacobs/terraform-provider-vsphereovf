// Code generated by counterfeiter. DO NOT EDIT.
package leasefakes

import (
	"sync"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/lease"
	"github.com/vmware/govmomi/nfc"
)

type FakeItemUpload struct {
	UploadStub        func(nfc.FileItem) error
	uploadMutex       sync.RWMutex
	uploadArgsForCall []struct {
		arg1 nfc.FileItem
	}
	uploadReturns struct {
		result1 error
	}
	uploadReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeItemUpload) Upload(arg1 nfc.FileItem) error {
	fake.uploadMutex.Lock()
	ret, specificReturn := fake.uploadReturnsOnCall[len(fake.uploadArgsForCall)]
	fake.uploadArgsForCall = append(fake.uploadArgsForCall, struct {
		arg1 nfc.FileItem
	}{arg1})
	fake.recordInvocation("Upload", []interface{}{arg1})
	fake.uploadMutex.Unlock()
	if fake.UploadStub != nil {
		return fake.UploadStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.uploadReturns.result1
}

func (fake *FakeItemUpload) UploadCallCount() int {
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	return len(fake.uploadArgsForCall)
}

func (fake *FakeItemUpload) UploadArgsForCall(i int) nfc.FileItem {
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	return fake.uploadArgsForCall[i].arg1
}

func (fake *FakeItemUpload) UploadReturns(result1 error) {
	fake.UploadStub = nil
	fake.uploadReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeItemUpload) UploadReturnsOnCall(i int, result1 error) {
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

func (fake *FakeItemUpload) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.uploadMutex.RLock()
	defer fake.uploadMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeItemUpload) recordInvocation(key string, args []interface{}) {
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

var _ lease.ItemUpload = new(FakeItemUpload)
