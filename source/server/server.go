// Copyright © 2022 Meroxa, Inc.
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

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
)

const (
	chanSize = 512
)

// HTTPServer server for listening post requests and convert it to sdk records.
type HTTPServer struct {
	server *http.Server
	port   int
	data   chan sdk.RawData
	errors chan error
}

func New(port int) *HTTPServer {
	data := make(chan sdk.RawData, chanSize)
	errors := make(chan error)

	return &HTTPServer{
		port:   port,
		data:   data,
		errors: errors,
	}
}

// Start http server.
func (h *HTTPServer) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.listen)

	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", h.port),
		Handler: mux,
	}

	go func() {
		err := h.server.ListenAndServe()
		if err != nil {
			h.errors <- err
		}
	}()
}

// Stop - shutdown channel and close channels.
func (h *HTTPServer) Stop(ctx context.Context) error {
	err := h.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	close(h.data)
	close(h.errors)

	return nil
}

func (h *HTTPServer) GetRecord(ctx context.Context) (sdk.Record, error) {
	select {
	case err := <-h.errors:
		return sdk.Record{}, fmt.Errorf("get record: %w", err)
	case data := <-h.data:
		pos, err := generatePosition()
		if err != nil {
			return sdk.Record{}, err
		}

		return sdk.Record{
			Position:  pos,
			CreatedAt: time.Now(),
			Payload:   data,
		}, nil
	case <-ctx.Done():
		return sdk.Record{}, ctx.Err()
	}
}

func (h *HTTPServer) listen(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusNotFound)

		return
	}

	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.errors <- fmt.Errorf("read body: %w", err)

		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	h.data <- rawBody

	w.WriteHeader(http.StatusCreated)

	return
}

func generatePosition() ([]byte, error) {
	return json.Marshal(uuid.New().String())
}
