// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.quantkernel.com
// Email       : ping.yang@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2025/05/06   youhei         New version
// -------------------------------------------------------------------

package exec

import (
	"bufio"
	"context"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/wengoldx/xcore/logger"
)

// Command executor.
//
// Use shell to execute any command scripts, set outHandler or errHandler
// to read console output (by 'echo' command print in shell script), or
// run command silent without any handlers sets.
//
//
// # NOTICE:
//	- The Executor enable execute system command like 'find', 'grep'...
//	- The command enable execute shell script file as './sample.sh argx'
//
//
// # USAGE:
//
//	// 1. Without any output handlers.
//	command := "./sample.sh arg1 arg2"
//	executor := exec.NewExecutor(command)
//	err := executor.Exec(context.Background())
//	// check err for execute success status.
//
//	// 2. With output handler.
//	command := "./sample.sh arg1 arg2"
//	executor := exec.NewExecutor(command,
//		exec.WithOutHandler(func(line string) {
//			// parse line string, and do samething...
//		}),
//	)
//	ctx, cancel := context.WithCancel(context.Background())
//	err := executor.Exec(ctx)
//	// check err for execute result, or call 'cancel' callback
//	// to cancel and stop executor all pipes.
//
//	// 3. User exec.WithOutHandler(), exec.WithErrHandler() to
//	// set both output and error handlers.
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
func NewExecutor(cmd string, opts ...Option) *Executor {
	executor := &Executor{command: cmd}
	for _, optfunc := range opts {
		optfunc(executor)
	}
	return executor
}

// Execute the given command and output logs from streaming hander.
func (ex *Executor) Exec(ctx context.Context) error {
	c := exec.CommandContext(ctx, "/bin/sh", "-c", ex.command)

	// set stdout pipe if exist read handler.
	readers := []*ConsoleReader{}
	if ex.outHandler != nil {
		if so, err := c.StdoutPipe(); err != nil {
			return err
		} else {
			r := &ConsoleReader{so, false}
			readers = append(readers, r)
		}
	}

	// set stderr pipe if exist read handler.
	if ex.errHandler != nil {
		if se, err := c.StderrPipe(); err != nil {
			return err
		} else {
			r := &ConsoleReader{se, true}
			readers = append(readers, r)
		}
	}

	// execute command and read console outputs.
	if pipes := len(readers); pipes > 0 {
		var wg sync.WaitGroup
		wg.Add(pipes)
		for _, reader := range readers {
			go ex.readOutputs(ctx, &wg, reader)
		}
		err := c.Start()
		wg.Wait()
		return err
	}
	return c.Run()
}

// Read output logs line by line and streaming by given handler function.
func (ex *Executor) readOutputs(ctx context.Context, wg *sync.WaitGroup, pipe *ConsoleReader) {
	// prepare pipe reader.
	reader := bufio.NewReader(pipe)
	defer wg.Done() // release lock when done.

	// read console output as line by line.
	for {
		select {
		case <-ctx.Done(): // interupt by cancel.
			return

		default: // read current line outputs.
			output, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
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
