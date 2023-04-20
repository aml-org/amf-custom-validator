module github.com/aml-org/amf-custom-validator

go 1.19

require (
	github.com/open-policy-agent/opa v0.47.0
	github.com/piprate/json-gold v0.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/flatbuffers v2.0.6+incompatible // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/yashtewari/glob-intersection v0.1.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/opencontainers/runc => github.com/opencontainers/runc v1.1.5

replace github.com/coreos/etcd => go.etcd.io/etcd/v3 v3.4.10

// replace github.com/dgrijalva/jwt-go => github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1

// replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.14.15

// replace github.com/pkg/sftp => github.com/pkg/sftp v1.11.0

// replace golang.org/x/image@0.5.0 => golang.org/x/image v0.5.0

// replace github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2

// exclude github.com/docker/docker v1.13.1