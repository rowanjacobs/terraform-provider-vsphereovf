// Code generated by counterfeiter. DO NOT EDIT.
package ovxfakes

import (
	"io"
	"sync"

	"github.com/rowanjacobs/terraform-provider-vsphereovf/vsphereovf/ovx"
)

type FakeReaderProvider struct {
	ReaderStub        func(string) (io.Reader, int64, error)
	readerMutex       sync.RWMutex
	readerArgsForCall []struct {
		arg1 string
	}
	readerReturns struct {
		result1 io.Reader
		result2 int64
		result3 error
	}
	readerReturnsOnCall map[int]struct {
		result1 io.Reader
		result2 int64
		result3 error
	}
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns     struct {
		result1 error
	}
	closeReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeReaderProvider) Reader(arg1 string) (io.Reader, int64, error) {
	fake.readerMutex.Lock()
	ret, specificReturn := fake.readerReturnsOnCall[len(fake.readerArgsForCall)]
	fake.readerArgsForCall = append(fake.readerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Reader", []interface{}{arg1})
	fake.readerMutex.Unlock()
	if fake.ReaderStub != nil {
		return fake.ReaderStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.readerReturns.result1, fake.readerReturns.result2, fake.readerReturns.result3
}

func (fake *FakeReaderProvider) ReaderCallCount() int {
	fake.readerMutex.RLock()
	defer fake.readerMutex.RUnlock()
	return len(fake.readerArgsForCall)
}

func (fake *FakeReaderProvider) ReaderArgsForCall(i int) string {
	fake.readerMutex.RLock()
	defer fake.readerMutex.RUnlock()
	return fake.readerArgsForCall[i].arg1
}

func (fake *FakeReaderProvider) ReaderReturns(result1 io.Reader, result2 int64, result3 error) {
	fake.ReaderStub = nil
	fake.readerReturns = struct {
		result1 io.Reader
		result2 int64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeReaderProvider) ReaderReturnsOnCall(i int, result1 io.Reader, result2 int64, result3 error) {
	fake.ReaderStub = nil
	if fake.readerReturnsOnCall == nil {
		fake.readerReturnsOnCall = make(map[int]struct {
			result1 io.Reader
			result2 int64
			result3 error
		})
	}
	fake.readerReturnsOnCall[i] = struct {
		result1 io.Reader
		result2 int64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeReaderProvider) Close() error {
	fake.closeMutex.Lock()
	ret, specificReturn := fake.closeReturnsOnCall[len(fake.closeArgsForCall)]
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.closeReturns.result1
}

func (fake *FakeReaderProvider) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeReaderProvider) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeReaderProvider) CloseReturnsOnCall(i int, result1 error) {
	fake.CloseStub = nil
	if fake.closeReturnsOnCall == nil {
		fake.closeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.closeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeReaderProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.readerMutex.RLock()
	defer fake.readerMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeReaderProvider) recordInvocation(key string, args []interface{}) {
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

var _ ovx.ReaderProvider = new(FakeReaderProvider)
