// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated from the elasticsearch-specification DO NOT EDIT.
// https://github.com/elastic/elasticsearch-specification/tree/4fcf747dfafc951e1dcf3077327e3dcee9107db3

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// SourceFilter type.
//
// https://github.com/elastic/elasticsearch-specification/blob/4fcf747dfafc951e1dcf3077327e3dcee9107db3/specification/_global/search/_types/SourceFilter.ts#L23-L31
type SourceFilter struct {
	Excludes []string `json:"excludes,omitempty"`
	Includes []string `json:"includes,omitempty"`
}

func (s *SourceFilter) UnmarshalJSON(data []byte) error {

	if !bytes.HasPrefix(data, []byte(`{`)) {
		var item string
		err := json.NewDecoder(bytes.NewReader(data)).Decode(&item)
		if err != nil {
			return err
		}
		s.Includes = append(s.Includes, item)
		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))

	for {
		t, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		switch t {

		case "excludes", "exclude":
			rawMsg := json.RawMessage{}
			dec.Decode(&rawMsg)
			if !bytes.HasPrefix(rawMsg, []byte("[")) {
				o := new(string)
				if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&o); err != nil {
					return fmt.Errorf("%s | %w", "Excludes", err)
				}

				s.Excludes = append(s.Excludes, *o)
			} else {
				if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&s.Excludes); err != nil {
					return fmt.Errorf("%s | %w", "Excludes", err)
				}
			}

		case "includes", "include":
			rawMsg := json.RawMessage{}
			dec.Decode(&rawMsg)
			if !bytes.HasPrefix(rawMsg, []byte("[")) {
				o := new(string)
				if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&o); err != nil {
					return fmt.Errorf("%s | %w", "Includes", err)
				}

				s.Includes = append(s.Includes, *o)
			} else {
				if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&s.Includes); err != nil {
					return fmt.Errorf("%s | %w", "Includes", err)
				}
			}

		}
	}
	return nil
}

// NewSourceFilter returns a SourceFilter.
func NewSourceFilter() *SourceFilter {
	r := &SourceFilter{}

	return r
}
