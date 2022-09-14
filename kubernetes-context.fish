#!/bin/env fish
kubectl config get-contexts | sed -i '1d'

kubectl config current-context