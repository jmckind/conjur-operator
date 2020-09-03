# Conjur OSS Operator

A Kubernetes Operator for managing Conjur OSS.

## Install

Start by cloning this repository and navigate to the root directory. Run the `deploy` target in the provided Makefile to deploy the operator. 

```
make deploy
```

## Usage

Once the operator is deployed, create a `Conjur` custom resource.

```
kubectl apply -f config/samples/conjur-basic.yaml
```

## License

The Conjur OSS Operator is released under the Apache 2.0 license. See the [LICENSE][license_file] file for details.

[license_file]:./LICENSE
