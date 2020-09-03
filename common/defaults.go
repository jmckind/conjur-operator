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
	// ConjurDefaultAuthenticators is the default authenticators to use for Conjur.
	ConjurDefaultAuthenticators = "authn"

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
	ConjurDefaultDatabaseVersion = "sha256:419a6454b74ff3b6fe5653757e107da3a26aee6cb34d625bbb7dac3637c12801" // 1-57

	// ConjurDefaultProxyImage is the default container image to use for the
	// Conjur server proxy.
	ConjurDefaultProxyImage = "nginx"

	// ConjurDefaultProxyProbePath is the default path to use for the Conjur
	// server proxy liveness and readiness probes.
	ConjurDefaultProxyProbePath = "/status"

	// ConjurDefaultProxySecurePort is the default secure port to use for the
	// Conjur server proxy.
	ConjurDefaultProxySecurePort = 9443

	// ConjurDefaultProxyUnsecurePort is the default unsecure port to use for
	// the Conjur server proxy.
	ConjurDefaultProxyUnsecurePort = 9000

	// ConjurDefaultProxyVersion is the default container image tag to use
	// for the Conjur server proxy.
	ConjurDefaultProxyVersion = "sha256:23b4dcdf0d34d4a129755fc6f52e1c6e23bb34ea011b315d87e193033bcd1b68" // 1.15

	// ConjurDefaultPVCCapicity is the default PVC capacity.
	ConjurDefaultPVCCapicity = "2Gi"

	// ConjurDefaultRSAKeySize is the default RSA key size when not specified.
	ConjurDefaultRSAKeySize = 2048

	// ConjurDefaultServerCount is the default replica count for the Conjur
	// server Deployment.
	ConjurDefaultServerCount = int32(1)

	// ConjurDefaultServerImage is the default container image to use for the
	// Conjur server.
	ConjurDefaultServerImage = "guygiat/conjur-oss-ocp"

	// ConjurDefaultServerPort is the default port to use for the Conjur server
	ConjurDefaultServerPort = 8000

	// ConjurDefaultServerProbePath is the default path to use for the Conjur
	// server liveness and readiness probes.
	ConjurDefaultServerProbePath = "/"

	// ConjurDefaultServerVersion is the default container image tag to use
	// for the Conjur server.
	ConjurDefaultServerVersion = "sha256:f2a4fa1dd75512223f1f58f12a8795756b19bf475ba8930915f4b46595e366ef" // Created 2020-07-12T16:15:02.2183617Z
)
