//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package view

import (
	"github.com/lastbackend/cli/pkg/util/table"
	"github.com/lastbackend/lastbackend/pkg/api/types/v1/views"
)

type NamespaceList []*Namespace
type Namespace views.Namespace

func (n *Namespace) Print() {

	println()
	table.PrintHorizontal(map[string]interface{}{
		"NAME":        n.Meta.Name,
		"DESCRIPTION": n.Meta.Description,
		"ENDPOINT":    n.Meta.Endpoint,
	})
	println()
}

func (nl *NamespaceList) Print() {

	t := table.New([]string{"NAME", "DESCRIPTION", "ENDPOINT"})
	t.VisibleHeader = true

	for _, n := range *nl {
		var data = map[string]interface{}{}

		data["NAME"] = n.Meta.Name
		data["DESCRIPTION"] = n.Meta.Description
		data["ENDPOINT"] = n.Meta.Endpoint

		t.AddRow(data)
	}

	println()
	t.Print()
	println()
}

func FromApiNamespaceView(namespace *views.Namespace) *Namespace {

	if namespace == nil {
		return nil
	}

	item := Namespace(*namespace)
	return &item
}

func FromApiNamespaceListView(namespaces *views.NamespaceList) *NamespaceList {
	var items = make(NamespaceList, 0)
	for _, namespace := range *namespaces {
		items = append(items, FromApiNamespaceView(namespace))
	}
	return &items
}
