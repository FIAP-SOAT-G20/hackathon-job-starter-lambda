#! /bin/bash

# Apply the configuration
kubectl apply -f fluent-bit-config.yaml

# Deploy Fluent Bit DaemonSet
kubectl apply -f https://raw.githubusercontent.com/fluent/fluent-bit-kubernetes-logging/master/fluent-bit-ds.yaml