// Copyright 2019 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package build

import (
	"context"
	"testing"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"golang.org/x/sync/errgroup"
)

type sleeper struct{}

var _ Interface = (*sleeper)(nil)

// IsSupportedReference implements Interface
func (r *sleeper) IsSupportedReference(ip string) bool {
	return true
}

// Build implements Interface
func (r *sleeper) Build(ip string) (v1.Image, error) {
	time.Sleep(50 * time.Millisecond)
	return nil, nil
}

func TestLimiter(t *testing.T) {
	b := NewLimiter(&sleeper{}, 2)

	start := time.Now()
	g, _ := errgroup.WithContext(context.TODO())
	for i := 0; i <= 10; i++ {
		g.Go(func() error {
			_, _ = b.Build("whatever")
			return nil
		})
	}
	g.Wait()

	// 50 ms * 10 builds / 2 concurrency = ~250ms
	if time.Now().Before(start.Add(250 * time.Millisecond)) {
		t.Fatal("Too many builds")
	}
}
