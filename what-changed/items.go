// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package what_changed

import (
	v2 "github.com/pb33f/libopenapi/datamodel/low/v2"
	v3 "github.com/pb33f/libopenapi/datamodel/low/v3"
)

type ItemsChanges struct {
	PropertyChanges
	ItemsChanges *ItemsChanges
}

func (i *ItemsChanges) TotalChanges() int {
	c := i.PropertyChanges.TotalChanges()
	if i.ItemsChanges != nil {
		c += i.ItemsChanges.TotalChanges()
	}
	return c
}

func (i *ItemsChanges) TotalBreakingChanges() int {
	c := i.PropertyChanges.TotalBreakingChanges()
	if i.ItemsChanges != nil {
		c += i.ItemsChanges.TotalBreakingChanges()
	}
	return c
}

func CompareItems(l, r *v2.Items) *ItemsChanges {

	var changes []*Change
	var props []*PropertyCheck

	ic := new(ItemsChanges)

	// header is identical to items, except for a description.
	props = append(props, addSwaggerHeaderProperties(l, r, &changes)...)
	CheckProperties(props)

	if !l.Items.IsEmpty() && !r.Items.IsEmpty() {
		// inline, check hashes, if they don't match, compare.
		if l.Items.Value.Hash() != r.Items.Value.Hash() {
			// compare.
			ic.ItemsChanges = CompareItems(l.Items.Value, r.Items.Value)
		}

	}
	if l.Items.IsEmpty() && !r.Items.IsEmpty() {
		// added items
		CreateChange(&changes, PropertyAdded, v3.ItemsLabel,
			nil, r.Items.GetValueNode(), true, nil, r.Items.GetValue())
	}
	if !l.Items.IsEmpty() && r.Items.IsEmpty() {
		// removed items
		CreateChange(&changes, PropertyRemoved, v3.ItemsLabel,
			l.Items.GetValueNode(), nil, true, l.Items.GetValue(),
			nil)
	}
	ic.Changes = changes
	if ic.TotalChanges() <= 0 {
		return nil
	}
	return ic
}