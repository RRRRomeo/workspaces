/*
 * @Author: hfouyang hfouyang@quant360.com
 * @Date: 2023-03-30 13:27:17
 * @LastEditors: hfouyang hfouyang@quant360.com
 * @LastEditTime: 2023-03-30 17:52:48
 * @FilePath: /githuab.com/RRRRomeo/workspaces/test_ldflags/test_flags.go
 */
package testldflags

import (
	"errors"
	"log"
)

var VERSION = "0.0.0.0"
var RELEASE = "0"

const TRUE = "1"

func Set(val string) error {
	if val == "" {
		return errors.New("val is invalid")
	}

	// ... sync dont support to async;
	VERSION = val
	return nil
}

func Get() string {
	if RELEASE != TRUE {
		log.Printf("start call get val func\n")
	}
	return VERSION
}
