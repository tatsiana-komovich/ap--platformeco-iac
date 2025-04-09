## About parameters
`enableIngressGateway` (default: `true`) - Should Istio use ingressgateways to access remote Pods?  
If remote Pods are accessible directly from our cluster (“flat” network), it is efficient not to use extra hop.  
`metadataEndpoint` (pattern: `^(https|file)://[0-9a-zA-Z._/-]+$`) - HTTPS endpoint with remote cluster metadata.