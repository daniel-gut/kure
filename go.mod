module github.com/daniel-gut/kure

go 1.14

// replace github.com/daniel-gut/kure/cmd => ./kure/cmd

// replace github.com/daniel-gut/kure/pkg/kure => ./kure/pkg/kure

require (
	github.com/daniel-gut/kure/cmd v0.0.0-00010101000000-000000000000
	github.com/daniel-gut/kure/pkg/kure v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/cobra v1.0.0
	k8s.io/api v0.18.1 // indirect
	k8s.io/client-go v11.0.0+incompatible // indirect
	k8s.io/utils v0.0.0-20200414100711-2df71ebbae66 // indirect
)
