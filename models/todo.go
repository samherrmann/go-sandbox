package models

import "fmt"

type ToDo []string

func (t *ToDo) Append(value string) {
	*t = append(*t, value)
}

func (t *ToDo) Remove(index int) error {
	if err := t.exists(index); err != nil {
		return err
	}
	*t = append((*t)[:index], (*t)[index+1:]...)
	return nil
}

func (t *ToDo) Update(index int, value string) error {
	if err := t.exists(index); err != nil {
		return err
	}
	(*t)[index] = value
	return nil
}

func (t *ToDo) exists(i int) error {
	if i >= len(*t) {
		return fmt.Errorf("index %v does not exist", i)
	}
	return nil
}
