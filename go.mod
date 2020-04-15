module github.com/daniel-gut/kure

go 1.14

replace github.com/daniel-gut/kure/cmd => ./kure/cmd

replace github.com/daniel-gut/kure/pkg/kure => ./kure/pkg/kure

replace github.com/daniel-gut/kure/pkg/clients => ./kure/pkg/clients

require (
	github.com/daniel-gut/kure/cmd v0.0.0-00010101000000-000000000000
	github.com/daniel-gut/kure/pkg/clients v0.0.0-00010101000000-000000000000 // indirect
	github.com/daniel-gut/kure/pkg/kure v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	k8s.io/api v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/utils v0.0.0-20191217112158-dcd0c905194b // indirect
	k8s.io/apimachinery v0.17.0
)
