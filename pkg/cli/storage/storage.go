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

package storage

import (
	"encoding/json"
	"github.com/lastbackend/cli/pkg/util/filesystem"
	"github.com/pkg/errors"
	"os"
	"strings"
)

var rootPath = strings.Join([]string{filesystem.HomeDir(), ".lastbackend"}, string(os.PathSeparator))
var fileToken = strings.Join([]string{rootPath, "token"}, string(os.PathSeparator))
var fileClusters = strings.Join([]string{rootPath, "clusters"}, string(os.PathSeparator))
var fileCluster = strings.Join([]string{rootPath, "cluster"}, string(os.PathSeparator))

func init() {
	os.MkdirAll(rootPath, 0700)
}

type Cluster struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
}

func SetToken(token string) error {
	return filesystem.WriteStrToFile(fileToken, token, 0700)
}

func GetToken() (string, error) {
	buf, err := filesystem.ReadFile(fileToken)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func AddLocalCluster(name, endpoint, token string, local bool) error {
	items, err := ListLocalCluster()
	if err != nil {
		return err
	}

	for _, item := range items {
		if item.Name == name {
			return errors.New("already exists")
		}
	}

	cluster := new(Cluster)
	cluster.Name = name
	cluster.Endpoint = endpoint
	cluster.Token = token

	items = append(items, cluster)

	buf, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return filesystem.WriteStrToFile(fileClusters, string(buf), 0700)
}

func ListLocalCluster() ([]*Cluster, error) {
	buf, err := filesystem.ReadFile(fileClusters)
	if err != nil {
		return nil, err
	}

	items := make([]*Cluster, 0)

	if len(buf) == 0 {
		return items, nil
	}

	if err := json.Unmarshal(buf, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func GetLocalCluster(name string) (*Cluster, error) {
	items, err := ListLocalCluster()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.Name == name {
			return item, nil
		}
	}

	return nil, nil
}

func DelLocalCluster(name string) error {

	items, err := ListLocalCluster()
	if err != nil {
		return err
	}

	for i, item := range items {
		if item.Name == name {

			cl, err := GetCluster()
			if err != nil {
				return err
			}

			match := strings.Split(cl, ".")
			if match[0] == "l" && match[1] == name {
				if err := SetCluster(""); err != nil {
					return err
				}
			}

			items = append(items[:i], items[i+1:]...)
			break
		}
	}

	buf, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return filesystem.WriteStrToFile(fileClusters, string(buf), 0700)
}

func SetCluster(cluster string) error {
	return filesystem.WriteStrToFile(fileCluster, cluster, 0700)
}

func GetCluster() (string, error) {

	buf, err := filesystem.ReadFile(fileCluster)
	if err != nil {
		return "", err
	}

	if len(buf) == 0 {
		return "", nil
	}

	return string(buf), nil
}
