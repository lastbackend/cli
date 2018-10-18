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

type ClusterList []*Cluster
type Cluster views.Cluster

func (cl *ClusterList) Print() {

	println()

	t := table.New([]string{"NAME"})
	t.VisibleHeader = true

	for _, s := range *cl {

		var data = map[string]interface{}{}
		data["NAME"] = s.Meta.Name
		t.AddRow(data)
	}
	println()
	t.Print()
	println()
}

func FromApiClusterListView(clusters *views.ClusterList) *ClusterList {

	var items = make(ClusterList, 0)
	for _, cluster := range *clusters {
		items = append(items, FromApiClusterView(cluster))
	}
	return &items
}

func (c *Cluster) Print() {

	println()
	table.PrintHorizontal(map[string]interface{}{
		"NAME":        c.Meta.Name,
		"DESCRIPTION": c.Meta.Description,
	})
	println()
}

func FromApiClusterView(cluster *views.Cluster) *Cluster {

	if cluster == nil {
		return nil
	}

	item := Cluster(*cluster)
	return &item
}
