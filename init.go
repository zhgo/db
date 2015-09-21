// Copyright 2014 The zhgo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
    "github.com/zhgo/config"
)

// Server
type conf struct {
    // Server
    DB  Server `json:"db"`

    // module list
    Modules map[string]confModule `json:"modules"`
}

// Module
type confModule struct {
    // DB Server
    DB  Server `json:"db"`

    // module name
    Name string `json:"name"`
}

func init() {
    var c conf

    // Load config file
    replaces := map[string]string{"{WorkingDir}": config.WorkingDir()}
    config.NewConfig("zhgo.json").Replace(replaces).Parse(c)

    // Module
    if c.Modules != nil {
        for k, v := range c.Modules {
            if v.Name == "" {
                v.Name = k
            }

            // db.Connections
            if v.DB.DSN == "" && c.DB.DSN != "" {
                v.DB = c.DB
            }

            if v.DB.DSN != "" {
                Servers[v.Name] = &v.DB
            }
        }
    }

}
