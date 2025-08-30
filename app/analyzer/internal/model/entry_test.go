package model

import (
	"reflect"
	"testing"
)

func TestEntryPlatformIteratorOrderByPriorityAsc(t *testing.T) {
	types := []EntryPlatformType{}
	for ep := range EntryPlatformIteratorOrderByPriorityAsc() {
		types = append(types, ep.Type())
	}

	expect := []EntryPlatformType{
		EntryPlatformTypeHatena,
		EntryPlatformTypeZenn,
	}
	if !reflect.DeepEqual(types, expect) {
		t.Errorf("expect %v, actual %v", expect, types)
	}
}
