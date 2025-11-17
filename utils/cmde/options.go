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

// Typed function to configure a Executor.
type Option func(*Executor)

// Set the out hanlder for Executor to read console normal outputs.
func WithOutHandler(out ConsoleHandler) Option {
	return func(ex *Executor) { ex.outHandler = out }
}

// Set the error hanlder for Executor to read console error outputs.
func WithErrHandler(err ConsoleHandler) Option {
	return func(ex *Executor) { ex.errHandler = err }
}

// Set the out and error hanlder for Executor to read all outputs.
func WithHandlers(cb ConsoleHandler) Option {
	return func(ex *Executor) { ex.outHandler, ex.errHandler = cb, cb }
}
