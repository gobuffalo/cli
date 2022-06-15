package test

import "context"

type Tester interface {
	Test(context.Context, string, []string) error
}
