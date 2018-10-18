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
	"os"
	"strings"

	"github.com/lastbackend/cli/pkg/util/filesystem"
)

var path = strings.Join([]string{filesystem.HomeDir(), ".lastbackend"}, string(os.PathSeparator))
var filepath = strings.Join([]string{path, "token"}, string(os.PathSeparator))

func SetToken(token string) error {
	err := os.MkdirAll(path, 755)
	if err != nil {
		return err
	}
	return filesystem.WriteStrToFile(filepath, token, 0644)
}

func GetToken() (string, error) {
	buf, err := filesystem.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
