package controller

import (
	"bfadmin/configuration"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

//
// Process supervisor, in its simplest form
// Needs to have proper synchronisation which is a TODO
//

type Runner struct {
	Name        string
	executable  string
	execParams  []string
	maxRunTime  time.Duration
	settingsDir string
	cmd         * exec.Cmd
	Status      string
}

func NewRunner(name string, executable string, settingsDir string, execParams []string, maxRunTimeMin int) *Runner {
	return &Runner{
		Name: name,
		executable:  executable,
		execParams: execParams,
		maxRunTime:  time.Duration(maxRunTimeMin) * time.Minute,
		settingsDir: settingsDir,
		Status:      "OFFLINE",
	}
}

func (r* Runner) Start() {
	if r.cmd != nil {
		log.Printf("Already started %s, PID=%d", r.executable, r.cmd.Process.Pid)
		return
	}
	r.initCmd()
	configuration.CopyConfigs(r.settingsDir)
	r.cmd.Start()
	if r.cmd.Process != nil {
		r.Status = "STARTING"
		log.Printf("Started %s, PID=%d", r.cmd, r.cmd.Process.Pid)
		go r.supervise()
	} else {
		log.Printf("Could not start %s", r.executable)
		r.cmd = nil
	}

}

func (r* Runner) initCmd() {
	r.cmd = exec.Command(r.executable, r.execParams ...)
	r.cmd.Dir = filepath.Dir(r.executable)
	r.cmd.Stdout = os.Stdout
}

func (r* Runner) supervise() {
	defer func() {
		r.cmd = nil
		r.Status = "OFFLINE"
	} ()
	for {
		shutdown := r.scheduleShutdown()
		err := r.cmd.Wait() // "waitpid"
		r.cancelShutdown(shutdown)

		if err == nil {
			log.Printf("Process PID=%d exited normally", r.cmd.Process.Pid)
			return
		}

		exitError, ok := err.(*exec.ExitError)
		if ok {
			log.Printf("Process PID=%d exited with code %d", r.cmd.Process.Pid, exitError.ExitCode())
			if exitError.ExitCode() == 1 || r.Status == "STOPPING" {
				return
			}
		} else {
			log.Printf("Unknown error during PID=%d shutdown, %T: %s", r.cmd.Process.Pid, err, err)
			return
		}

		r.initCmd()
		r.cmd.Start()
		r.Status = "STARTING"
		log.Printf("Restarted %s, PID=%d", r.executable, r.cmd.Process.Pid)
	}
}

func (r* Runner) scheduleShutdown() chan<- bool {
	if r.maxRunTime > 0 {
		shutdown := make(chan bool)
		timer := time.NewTimer(r.maxRunTime)
		go func() {
			log.Printf("Setting up shutdown timer for %s", r.Name)
			select {
			case <-timer.C:
				log.Printf("Max running time eached for %s", r.Name)
				r.Stop()
				<-shutdown // clean up shutdown chan
			case <-shutdown:
				if !timer.Stop() { // clean up chan (only if timer fired after we received shutdown)
					<-timer.C
				}
				log.Printf("Shutdown timer stopped for %s", r.Name)
			}
		}()
		return shutdown
	}
	return nil
}

func (r Runner) cancelShutdown(shutdown chan<- bool) {
	if shutdown != nil {
		shutdown <- true
	}
}

func (r* Runner) Stop() {
	if r.cmd != nil {
		log.Printf("Stopping PID=%d", r.cmd.Process.Pid)
		r.Status = "STOPPING"
		if r.Name == "bf2" || r.Name == "bf2142" { // BF2/BF2142 are special =(
			r.cmd.Process.Signal(syscall.Signal(15))
		} else {
			r.cmd.Process.Signal(syscall.Signal(2))
		}
	}
}