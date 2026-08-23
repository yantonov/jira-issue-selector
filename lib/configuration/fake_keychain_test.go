package configuration

import "fmt"

type FakeKeychain struct {
	Values  map[string]string
	Failure error
	Reads   []string
}

func NewFakeKeychain() *FakeKeychain {
	return &FakeKeychain{Values: map[string]string{}}
}

func (e *FakeKeychain) Get(key string) (string, error) {
	e.Reads = append(e.Reads, key)
	if e.Failure != nil {
		return "", e.Failure
	}
	return e.Values[key], nil
}

func (e *FakeKeychain) Set(key string, value string) error {
	if e.Failure != nil {
		return e.Failure
	}
	e.Values[key] = value
	return nil
}

func (e *FakeKeychain) Delete(key string) error {
	if e.Failure != nil {
		return e.Failure
	}
	delete(e.Values, key)
	return nil
}

func (e *FakeKeychain) String() string {
	return fmt.Sprintf("%v", e.Values)
}
