// Copyright 2019 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

type Config struct {
	Stocks *StocksConfig `yaml:"stocks"`
}

type StocksConfig struct {
	Symbols []string `yaml:"symbols"`
}
