package tests

import (
	"os"
	"sync"
	"testing"
	"time"

	ws "github.com/ALiwoto/ssg/ssg"
	"github.com/ALiwoto/ssg/ssg/shellUtils"
)

func TestShell01(t *testing.T) {
	result := ws.RunCommand("go version")
	if result == nil {
		t.Error("result is nil")
		return
	}

	if result.Error != nil {
		t.Error(result.Error)
		return
	}

	if result.Stdout == "" {
		t.Error("unexpected empty stdout string from result")
		return
	}
}

func TestShellAsync01(t *testing.T) {
	if os.PathSeparator != '/' {
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	result := ws.RunCommandAsync("sleep 100 && echo hello")
	if result == nil {
		t.Error("result is nil")
		return
	}

	result.WaitAndRun(time.Second, 3*time.Second, func(r *ws.ExecuteCommandResult) {
		if r.IsDone() {
			t.Error("unexpected true value returned from IsDone method")
			return
		}

		// kill the process
		err := r.Kill()
		if err != nil {
			t.Error("when tried to kill the proess: ", err)
			return
		}

		//log.Println("killed the proccess")
		//time.Sleep(300 * time.Second)

		wg.Done()
	})

	wg.Wait()
}

func TestShellAsyncKillIsRaceFree(t *testing.T) {
	finished := make(chan bool, 1)
	result := shellUtils.ExecuteCommandAsync("--", &shellUtils.ExecuteCommandConfig{
		TargetRunner:  os.Args[0],
		PrimaryArgs:   []string{"-test.run=^TestShellAsyncHelperProcess$"},
		AdditionalEnv: append(os.Environ(), "SSG_ASYNC_HELPER_PROCESS=1"),
		FinishedChan:  finished,
		IsAsync:       true,
	})
	if result == nil {
		t.Fatal("result is nil")
	}

	killResult := make(chan error, 1)
	result.WaitAndRun(5*time.Millisecond, 25*time.Millisecond, func(result *shellUtils.ExecuteCommandResult) {
		killResult <- result.Kill()
	})

	select {
	case err := <-killResult:
		if err != nil {
			t.Fatalf("failed to kill async command: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting to kill async command")
	}

	select {
	case <-finished:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for killed command to finish")
	}
}

func TestShellAsyncHelperProcess(t *testing.T) {
	if os.Getenv("SSG_ASYNC_HELPER_PROCESS") != "1" {
		return
	}

	for {
		time.Sleep(time.Second)
	}
}
