// Copyright 2020 The Conjur Operator Developers

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"time"
)

const (
	// ConjurAppName is the application name that is used in labels and annotations.
	ConjurAppName = "conjur-oss"

	// ConjurComponentPostgres is the value to use for the postgres component.
	ConjurComponentPostgres = "postgres"

	// ConjurDuration365Days is a duration representing 365 days.
	ConjurDuration365Days = time.Hour * 24 * 365

	// Version is the operator version.
	Version = "0.0.1"
)
