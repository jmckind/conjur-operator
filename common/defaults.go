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

const (
	// ConjurDefaultTLSFileMode is the default file mode to use when mounting
	// TLS secrets.
	ConjurDefaultTLSFileMode = int32(384)

	// ConjurDefaultDatabaseCount is the default replica count for the Conjur
	// database StatefulSet.
	ConjurDefaultDatabaseCount = int32(1)

	// ConjurDefaultDatabaseImage is the default container image to use for the
	// Conjur database.
	ConjurDefaultDatabaseImage = "registry.access.redhat.com/rhscl/postgresql-10-rhel7"

	// ConjurDefaultDatabasePort is the default port to use for the Conjur
	// database.
	ConjurDefaultDatabasePort = 5432

	// ConjurDefaultDatabaseVersion is the default container image tag to use
	// for the Conjur database.
	ConjurDefaultDatabaseVersion = "sha256:419a6454b74ff3b6fe5653757e107da3a26aee6cb34d625bbb7dac3637c12801"

	// ConjurDefaultPVCCapicity is the default PVC capacity.
	ConjurDefaultPVCCapicity = "2Gi"

	// ConjurDefaultRSAKeySize is the default RSA key size when not specified.
	ConjurDefaultRSAKeySize = 2048
)
