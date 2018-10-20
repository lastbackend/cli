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
	"github.com/pkg/errors"
	"os"
	"strings"

	"github.com/lastbackend/cli/pkg/util/filesystem"
)

var rootPath = strings.Join([]string{filesystem.HomeDir(), ".lastbackend"}, string(os.PathSeparator))
var fileToken = strings.Join([]string{rootPath, "token"}, string(os.PathSeparator))
var fileClusters = strings.Join([]string{rootPath, "clusters"}, string(os.PathSeparator))
var fileCluster = strings.Join([]string{rootPath, "cluster"}, string(os.PathSeparator))

func init() {
	os.MkdirAll(rootPath, 0700)
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

func AddLocalCluster(name, endpoint string) error {
	items, err := ListLocalCluster()
	if err != nil {
		return err
	}

	if _, ok := items[name]; ok {
		return errors.New("already exists")
	}

	items[name] = endpoint

	buf, err := json.Marshal(items)
	if err != nil {
		return err
	}

	return filesystem.WriteStrToFile(fileClusters, string(buf), 0700)
}

func ListLocalCluster() (map[string]string, error) {
	buf, err := filesystem.ReadFile(fileClusters)
	if err != nil {
		return nil, err
	}

	items := make(map[string]string, 0)

	if len(buf) == 0 {
		return items, nil
	}

	if err := json.Unmarshal(buf, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func GetLocalCluster(name string) (string, error) {
	items, err := ListLocalCluster()
	if err != nil {
		return "", err
	}

	if endpoint, ok := items[name]; ok {
		return endpoint, nil
	}

	return "", nil
}

func DelLocalCluster(name string) error {

	items, err := ListLocalCluster()
	if err != nil {
		return err
	}

	if _, ok := items[name]; ok {
		delete(items, name)

		buf, err := json.Marshal(items)
		if err != nil {
			return err
		}

		return filesystem.WriteStrToFile(fileClusters, string(buf), 0700)
	}

	return nil
}

type Cluster struct {
	Local    bool   `json:"local"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint,omitempty"`
}

func SetCluster(name, endpoint string, local bool) error {

	var c = new(Cluster)
	c.Name = name
	c.Local = local
	c.Endpoint = endpoint

	buf, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return filesystem.WriteStrToFile(fileCluster, string(buf), 0700)
}

func GetCluster() (*Cluster, error) {

	var res = new(Cluster)

	buf, err := filesystem.ReadFile(fileCluster)
	if err != nil {
		return nil, err
	}

	if len(buf) == 0 {
		return nil, nil
	}

	if err := json.Unmarshal(buf, &res); err != nil {
		return nil, err
	}

	return nil, nil
}
