// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package secure

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/wengoldx/xcore/invar"
)

// Random coder to generate unique number code
//
// `USEAGE` :
//
//	coder := mvc.NewSoleCoder()
//	code, _ := coder.Gen(6)
//	logger.I("6 chars code:", code)
//
//	code, _ = coder.Gen(8)
//	logger.I("8 chars code:", code)
//
//	code, _ := coder.Gen(6, 5)
//	logger.I("max retry 5 times, 6 chars code:", code)
//
//	code, _ = coder.Gen(8, 5)
//	logger.I("max retry 5 times, 8 chars code:", code)
type SoleCoder struct {
	codes map[string]bool
}

// Create SoleCoder and init with exist codes
func NewSoleCoder(data ...[]string) *SoleCoder {
	coder := &SoleCoder{
		codes: make(map[string]bool),
	}

	datalen := len(data)
	if datalen > 0 {
		for i := 0; i < datalen; i++ {
			if data[i] != nil && len(data[i]) > 0 {
				for _, code := range data[i] {
					coder.codes[code] = true
				}
			}
		}
	}
	return coder
}

// Generate a given length number code, it may throw a error
// when over the retry times
func (c *SoleCoder) Gen(codelen int, times ...int) (string, error) {
	if c.codes == nil {
		c.codes = make(map[string]bool)
	}

	radix := 1
	for i := 0; i < codelen; i++ {
		radix *= 10
	}

	if len(times) > 0 {
		for i := 0; i < times[0]; i++ {
			code, err := c.innerGenerate(codelen, radix)
			if err != nil {
				continue
			}
			return code, nil
		}
		return "", invar.ErrOverTimes
	}
	return c.innerGenerate(codelen, radix)
}

// Remove used sole code outof cache
func (c *SoleCoder) Remove(code string) {
	if c.codes != nil {
		if _, ok := c.codes[code]; ok {
			c.codes[code] = false
		}
	}
}

// Generate a give length number code
func (c *SoleCoder) innerGenerate(codelen, radix int) (string, error) {
	rander.Seed(time.Now().UnixNano())
	format := "%0" + fmt.Sprintf("%d", codelen) + "d"
	code := fmt.Sprintf(format, rand.Intn(radix))

	// check generated code if it unique
	if cc, ok := c.codes[code]; ok && cc {
		return "", invar.ErrDupData
	}
	c.codes[code] = true
	return code, nil
}
