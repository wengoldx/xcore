// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package mvc

type Options struct {
	datatype string // Response data type, default 'json', set one of 'json', 'jsonp', 'xml', 'yaml'.
	validate bool   // Need validate parsed input params.
	datas    any    // Response datas for Respon() function options.

	Protect bool // Response the error messages to frontend with 40x status codes.
	Silent  bool // Stop output debug logs when API controller executing.
}

type Option func(*Options)

func WithDataType(datatype string) Option {
	return func(opts *Options) {
		opts.datatype = datatype
	}
}

func WithDatas(datas any) Option {
	return func(opts *Options) {
		opts.datas = datas
	}
}

func WithProtect(protect bool) Option {
	return func(opts *Options) {
		opts.Protect = protect
	}
}

func WithSlie(silent bool) Option {
	return func(opts *Options) {
		opts.Silent = silent
	}
}

func newOptions(protect, silent bool) *Options {
	return &Options{
		datatype: "json",
		Protect:  protect, Silent: silent,
	}
}

func parseOptions(validate bool, options ...Option) *Options {
	opts := &Options{validate: validate, Protect: true}
	for _, optFunc := range options {
		optFunc(opts)
	}

	if opts.datatype == "" {
		opts.datatype = "json"
	}
	return opts
}

func (opts *Options) outType(datatype string) *Options {
	opts.datatype = datatype
	return opts
}
