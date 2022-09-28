package model

import (
	"errors"
	"fmt"
)

func run(data Handler, name string, isQuick bool, opts ...opt) error {
	if modelIsFull() {
		return errors.New("model is full")
	}
	if name == NullString {
		return errors.New("model name is required")
	}

	m := newModel(data, name, opts...)
	err := registerModel(m)
	if err != nil {
		return err
	}
	go m.wait(isQuick)
	fmt.Println(name, "start...")
	return nil
}

func modelIsFull() bool {
	return modelCount() > MaxModelCount
}
