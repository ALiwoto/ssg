package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/ALiwoto/ssg/ssg/mapUtils"
)

func TestSafeEMapGetWithOptionsPreventsMutexLifetimeRace(t *testing.T) {
	const key = "resource"

	m := mapUtils.NewSafeEMap[string, sync.Mutex]()
	m.SetExpiration(time.Hour)

	original := &sync.Mutex{}
	original.Lock()
	m.Add(key, original)

	doFnStarted := make(chan struct{})
	valueLocked := make(chan struct{})
	getReturned := make(chan struct{})
	releaseWorker := make(chan struct{})
	workerDone := make(chan struct{})

	var workerValue *sync.Mutex
	createCalled := false
	go func() {
		defer close(workerDone)

		workerValue = m.GetWithOptions(key, &mapUtils.GetOptions[string, sync.Mutex]{
			CreateFn: func() (*sync.Mutex, bool) {
				createCalled = true
				return &sync.Mutex{}, true
			},
			DoFn: func(value *sync.Mutex) {
				close(doFnStarted)
				value.Lock()
				close(valueLocked)
			},
		})
		close(getReturned)

		<-releaseWorker
		workerValue.Unlock()
	}()

	select {
	case <-doFnStarted:
	case <-time.After(time.Second):
		t.Fatal("DoFn did not start")
	}

	cleanupStarted := make(chan struct{})
	cleanupDone := make(chan struct{})
	deleted := false
	go func() {
		close(cleanupStarted)
		m.DeleteIf(key, func(value *sync.Mutex) bool {
			if !value.TryLock() {
				return false
			}

			value.Unlock()
			deleted = true
			return true
		})
		close(cleanupDone)
	}()

	<-cleanupStarted
	select {
	case <-cleanupDone:
		t.Fatal("cleanup passed GetWithOptions while DoFn was waiting for the value lock")
	case <-time.After(50 * time.Millisecond):
	}

	original.Unlock()

	select {
	case <-valueLocked:
	case <-time.After(time.Second):
		t.Fatal("DoFn did not acquire the value lock")
	}

	select {
	case <-cleanupDone:
	case <-time.After(time.Second):
		t.Fatal("cleanup remained blocked after GetWithOptions returned")
	}
	select {
	case <-getReturned:
	case <-time.After(time.Second):
		t.Fatal("GetWithOptions did not return after DoFn acquired the value lock")
	}

	if createCalled {
		t.Fatal("CreateFn was called for an existing mutex")
	}
	if workerValue != original {
		t.Fatal("GetWithOptions returned a different mutex")
	}
	if deleted || !m.Exists(key) {
		t.Fatal("cleanup deleted the mutex while it was locked by the worker")
	}

	close(releaseWorker)
	select {
	case <-workerDone:
	case <-time.After(time.Second):
		t.Fatal("worker did not release the value lock")
	}

	m.DeleteIf(key, func(value *sync.Mutex) bool {
		if !value.TryLock() {
			return false
		}

		value.Unlock()
		return true
	})
	if m.Exists(key) {
		t.Fatal("cleanup did not delete the unused mutex")
	}

	createCalled = false
	replacement := m.GetWithOptions(key, &mapUtils.GetOptions[string, sync.Mutex]{
		CreateFn: func() (*sync.Mutex, bool) {
			createCalled = true
			return &sync.Mutex{}, true
		},
		DoFn: func(value *sync.Mutex) {
			value.Lock()
		},
	})
	if replacement == nil {
		t.Fatal("GetWithOptions did not create a replacement mutex")
	}
	replacement.Unlock()

	if !createCalled || replacement == original {
		t.Fatal("GetWithOptions did not create a new mutex after safe cleanup")
	}
}
