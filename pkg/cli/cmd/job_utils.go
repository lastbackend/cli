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

package cmd

import (
	"fmt"
	"github.com/lastbackend/lastbackend/pkg/distribution/errors"
	"github.com/lastbackend/lastbackend/pkg/distribution/types"
	"strings"
)

func jobParseSelfLink(selflink string) (string, string, error) {
	match := strings.Split(selflink, "/")

	var (
		namespace, name string
	)

	switch len(match) {
	case 2:
		namespace = match[0]
		name = match[1]
	case 1:
		fmt.Println("Use default namespace:", types.DEFAULT_NAMESPACE)
		namespace = types.DEFAULT_NAMESPACE
		name = match[0]
	default:
		return "", "", errors.New("invalid service name provided")
	}

	return namespace, name, nil
}
