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

package putdatalifecycle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Request holds the request body struct for the package putdatalifecycle
//
// https://github.com/elastic/elasticsearch-specification/blob/4fcf747dfafc951e1dcf3077327e3dcee9107db3/specification/indices/put_data_lifecycle/IndicesPutDataLifecycleRequest.ts#L25-L76
type Request struct {

	// DataRetention If defined, every document added to this data stream will be stored at least
	// for this time frame.
	// Any time after this duration the document could be deleted.
	// When empty, every document in this data stream will be stored indefinitely.
	DataRetention types.Duration `json:"data_retention,omitempty"`
	// Downsampling If defined, every backing index will execute the configured downsampling
	// configuration after the backing
	// index is not the data stream write index anymore.
	Downsampling *types.DataStreamLifecycleDownsampling `json:"downsampling,omitempty"`
}

// NewRequest returns a Request
func NewRequest() *Request {
	r := &Request{}

	return r
}

// FromJSON allows to load an arbitrary json into the request structure
func (r *Request) FromJSON(data string) (*Request, error) {
	var req Request
	err := json.Unmarshal([]byte(data), &req)

	if err != nil {
		return nil, fmt.Errorf("could not deserialise json into Putdatalifecycle request: %w", err)
	}

	return &req, nil
}

func (s *Request) UnmarshalJSON(data []byte) error {
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

		case "data_retention":
			if err := dec.Decode(&s.DataRetention); err != nil {
				return fmt.Errorf("%s | %w", "DataRetention", err)
			}

		case "downsampling":
			if err := dec.Decode(&s.Downsampling); err != nil {
				return fmt.Errorf("%s | %w", "Downsampling", err)
			}

		}
	}
	return nil
}
