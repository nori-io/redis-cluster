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
	"testing"
	"time"

	"github.com/cheebo/go-config"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"

	"github.com/nori-io/nori/core/plugins/mocks"
)

const (
	testRedisAddr = "localhost:6379"
)

var (
	testKey   = "testkey"
	testValue = []byte("testvalue")
)

func TestPackage(t *testing.T) {
	assert := assert.New(t)

	registry := new(mocks.Registry)

	cfg := go_config.New()
	cfg.SetDefault("redis.cluster_addrs", testRedisAddr)

	registry.On("Config").Return(cfg)

	p := new(plugin)

	assert.NotNil(p.Meta())
	assert.NotEmpty(p.Meta().GetDescription().Name)

	p.Start(nil, registry)

	redis, ok := p.Instance().(redis.UniversalClient)
	assert.True(ok)
	assert.NotNil(redis)

	err := redis.Set(testKey, testValue, time.Duration(0)).Err()
	assert.Nil(err)

	bs, err := redis.Get(testKey).Bytes()
	assert.Nil(err)
	assert.Equal(testValue, bs)

	err = redis.Del(testKey).Err()
	assert.Nil(err)

	assert.Nil(p.Stop(nil, nil))
}
