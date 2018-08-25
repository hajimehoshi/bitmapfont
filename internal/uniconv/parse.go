// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package uniconv

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func Parse(f io.Reader, split string) (map[int]rune, error) {
	m := map[int]rune{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		tokens := strings.Split(line, split)
		if len(tokens) != 2 {
			continue
		}
		from, err := strconv.ParseInt(tokens[0], 0, 32)
		if err != nil {
			return nil, err
		}
		uni, err := strconv.ParseInt(tokens[1], 0, 32)
		if err != nil {
			return nil, err
		}
		m[int(from)] = rune(uni)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return m, nil
}
