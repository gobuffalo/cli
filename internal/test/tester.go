package test

import "context"

type Tester interface {
	Name() string
	Test(context.Context, string, []string) error
}
