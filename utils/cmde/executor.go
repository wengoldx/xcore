// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.quantkernel.com
// Email       : ping.yang@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/05/06   youhei         New version
// -------------------------------------------------------------------

package cmde

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"

	"github.com/wengoldx/xcore/logger"
)

// Command executor.
//
// Use shell to execute any command scripts, set outHandler or errHandler
// to read console output (by 'echo' command print in shell script), or
// run command silent without any handlers sets.
//
// # NOTICE:
//	- The Executor enable execute system command like 'find', 'grep'...
//	- The command enable execute shell script file as './sample.sh argx'
type Executor struct {
	command    string
	outHandler ConsoleHandler
	errHandler ConsoleHandler
}

// Callback handler for read console output line by line.
//
// # WARING:
//	- The 'line' string will tirm '\n' end of line.
type ConsoleHandler func(line string)

// Read console outputs as string.
type ConsoleReader struct {
	io.ReadCloser      // Console outputs reader.
	isStderr      bool // Indicate this reader whetcher stdout or stderr.
}

// Create a command executor to execute command.
//
// # USAGE:
//
// 1. Without any output handlers.
//
//	command := "./sample.sh arg1 arg2"
//	executor := cmd.NewExecutor(command)
//	err := executor.Exec(context.Background())
//	// check err for execute success status.
//
// 2. With output handler.
//
//	command := "./sample.sh arg1 arg2"
//	executor := cmd.NewExecutor(command,
//		cmd.WithOutHandler(func(line string) {
//			// parse line string, and do samething...
//		}),
//	)
//	ctx, cancel := context.WithCancel(context.Background())
//	err := executor.Exec(ctx)
//	// check err for execute result, or call 'cancel' callback
//	// to cancel and stop executor all pipes.
//
// 3. Async execute command.
//
//	done := make(chan error, 1)
//	command := "./sample.sh arg1 arg2"
//	executor := cmd.NewExecutor(command) // enable set handlers.
//	err := executor.Async(context.Background(), done)
//	// check err, do anythins here!
//	// ...
//	err = <-done // Wait command finished!
//
// 4. User cmd.WithOutHandler(), cmd.WithErrHandler() to
// set both output and error handlers.
func NewExecutor(cmd string, opts ...Option) *Executor {
	executor := &Executor{command: cmd}
	for _, optfunc := range opts {
		optfunc(executor)
	}
	return executor
}

// Execute the command and output logs from pipe handers.
//
// # NOTICE:
//
// This method will sync execute command until finished.
//
//	See cmde.Async() to execute command on async way.
func (ex *Executor) Exec(ctx context.Context) error {
	c := exec.CommandContext(ctx, "/bin/sh", "-c", ex.command)
	readers, err := ex.setupReaders(c)
	if err != nil || readers == nil {
		return err
	}

	// execute and read console outputs.
	close := make(chan struct{}, 1)
	for _, reader := range readers {
		go ex.readOutputs(ctx, reader, close)
	}
	err = c.Run()       // wait finished.
	close <- struct{}{} // close pipes.
	return err
}

// Async execute the command and output logs from pipe handers.
//
// # WARNING:
//
// This method will async execute command not wait it finished, so request
// the caller must blocking to wait 'notify' return if set any handlers.
//
//	See cmde.Exec() to execute command on sync way.
func (ex *Executor) Async(ctx context.Context, notify chan error) error {
	c := exec.CommandContext(ctx, "/bin/sh", "-c", ex.command)
	if notify != nil {
		readers, err := ex.setupReaders(c)
		if err != nil || readers == nil {
			return err
		}

		// execute and read console outputs.
		close := make(chan struct{}, 1)
		for _, reader := range readers {
			go ex.readOutputs(ctx, reader, close)
		}

		err = c.Start()
		go func() {
			notify <- c.Wait()  // wait finished.
			close <- struct{}{} // close pipes.
		}()
		return err
	}
	return c.Start() // not wait finished!
}

// Set stdout or stderr pipe outputs readers when user set the handler options.
func (ex *Executor) setupReaders(cmd *exec.Cmd) ([]*ConsoleReader, error) {
	// set stdout pipe if exist read handler.
	readers := []*ConsoleReader{}
	if ex.outHandler != nil {
		if so, err := cmd.StdoutPipe(); err != nil {
			return nil, err
		} else {
			r := &ConsoleReader{so, false}
			readers = append(readers, r)
		}
	}

	// set stderr pipe if exist read handler.
	if ex.errHandler != nil {
		if se, err := cmd.StderrPipe(); err != nil {
			return nil, err
		} else {
			r := &ConsoleReader{se, true}
			readers = append(readers, r)
		}
	}
	return readers, nil
}

// Read output logs line by line and streaming by given handler function.
func (ex *Executor) readOutputs(ctx context.Context, pipe *ConsoleReader, close chan struct{}) {
	// prepare pipe reader.
	reader := bufio.NewReader(pipe)

	// read console output as line by line.
	for {
		select {
		case <-ctx.Done(): // interupt by cancel.
			return
		case <-close: // close when command finished.
			logger.D("Exist pipe reader.")
			return

		default: // read current line outputs.
			output, err := reader.ReadString('\n')
			if err != nil {
				isclosed := strings.Contains(err.Error(), "file already closed")
				if err != io.EOF && !isclosed {
					logger.E("Read line, err:", err)
				}
				return
			}
			output = strings.TrimSuffix(output, "\n")
			if pipe.isStderr {
				ex.errHandler(output)
			} else {
				ex.outHandler(output)
			}
		}
	}
}
