package clio

import "fmt"

// StringSliceVar is a flag that takes a string slice value
// will be used across the cli to receive a list of strings
type StringSliceVar []string

func (i *StringSliceVar) String() string {
	return fmt.Sprintf("%s", *i)
}

func (i *StringSliceVar) Set(value string) error {
	*i = append(*i, value)

	return nil
}
