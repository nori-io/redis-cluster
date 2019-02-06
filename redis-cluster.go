// Copyright (C) 2018 The Nori Authors info@nori.io
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program; if not, see <http://www.gnu.org/licenses/>.
package main

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-redis/redis"

	cfg "github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori-common/meta"
	noriPlugin "github.com/nori-io/nori-common/plugin"
)

type plugin struct {
	instance   redis.UniversalClient
	pingEnable bool
	config     config
}

type config struct {
	clusterAddr cfg.String
}

var (
	Plugin plugin
)

const (
	pingTimeout      time.Duration = 10 * time.Second
	defaultSeparator               = " "
)

func (p *plugin) Init(_ context.Context, configManager cfg.Manager) error {
	cm := configManager.Register(p.Meta())
	p.config = config{
		clusterAddr: cm.String("redis.cluster_addrs", ""),
	}
	return nil
}

func (p *plugin) Instance() interface{} {
	return p.instance
}

func (p plugin) Meta() meta.Meta {
	return &meta.Data{
		ID: meta.ID{
			ID:      "nori/redis/cluster",
			Version: "1.0.0",
		},
		Author: meta.Author{
			Name: "Nori",
			URI:  "https://nori.io",
		},
		Core: meta.Core{
			VersionConstraint: ">=1.0.0, <2.0.0",
		},
		Dependencies: []meta.Dependency{},
		Description: meta.Description{
			Name: "Nori: Redis Cluster",
		},
		Interface: meta.Custom,
		License: meta.License{
			Title: "",
			Type:  "GPLv3",
			URI:   "https://www.gnu.org/licenses/"},
		Tags: []string{"cache", "redis"},
	}

}

func (p *plugin) Start(ctx context.Context, registry noriPlugin.Registry) error {
	if p.instance == nil {
		addrs := strings.Split(p.config.clusterAddr(), defaultSeparator)

		if len(addrs) == 0 {
			return errors.New("Redis address empty")
		}

		p.instance = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs: addrs,
		})

		go func(instance redis.UniversalClient) {
			var err error
			//logger := registry.GetLogger()
			for {
				err = instance.Ping().Err()
				if err != nil {
					// logger.Error(err)
					err = instance.Close()
					// logger.Error(err)
					instance = redis.NewUniversalClient(&redis.UniversalOptions{
						Addrs: addrs,
					})
				}

				time.Sleep(pingTimeout)

			}
		}(p.instance)
	}
	return nil
}

func (p *plugin) Stop(_ context.Context, _ noriPlugin.Registry) error {
	err := p.instance.Close()
	p.instance = nil
	return err
}
