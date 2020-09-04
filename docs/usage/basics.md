# Basic Usage

See the [Conjur Reference][conjur_reference] for the full list of properties and defaults to configure the Conjur environment.

The following example shows the most minimal valid manifest to create a new Conjur environment with the default configuration.

```yaml
apiVersion: oss.cyberark.com/v1alpha1
kind: Conjur
metadata:
  name: example
spec: {}
```

[conjur_reference]:../reference/conjur.md
