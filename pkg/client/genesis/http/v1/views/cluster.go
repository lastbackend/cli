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

package views

import "github.com/lastbackend/lastbackend/pkg/api/types/v1/views"

// swagger:ignore
// ClusterList is a list of cluster models for api
//
// swagger:model views_cluster_list
type ClusterList []*ClusterView

type ClusterView struct {
	Meta   views.Meta          `json:"meta"`
	Status views.ClusterStatus `json:"status"`
	Spec   ClusterSpec         `json:"spec"`
}

type ClusterSpec struct {
	Endpoint string `json:"endpoint"`
}
