module github.com/daniel-gut/kure

go 1.14

replace github.com/daniel-gut/kure/cmd => ./kure/cmd

replace github.com/daniel-gut/kure/pkg/kure => ./kure/pkg/kure

replace github.com/daniel-gut/kure/pkg/clients => ./kure/pkg/clients

require (
	github.com/araddon/dateparse v0.0.0-20200409225146-d820a6159ab1 // indirect
	github.com/aybabtme/uniplot v0.0.0-20151203143629-039c559e5e7e // indirect
	github.com/daniel-gut/kure/cmd v0.0.0-00010101000000-000000000000
	github.com/daniel-gut/kure/pkg/clients v0.0.0-00010101000000-000000000000 // indirect
	github.com/daniel-gut/kure/pkg/kure v0.0.0-00010101000000-000000000000 // indirect
	github.com/spf13/cobra v1.0.0 // indirect
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/utils v0.0.0-20191217112158-dcd0c905194b // indirect
)
