package main

import (
	"golang.org/x/net/context"
)

type Procedure func(context.Context, interface{}) (interface{}, error)
