@Echo.
@cd internal
@pwd
@Echo.
go generate -run "genny.*" .
go generate -run "bundle.*" .

@Echo.
@cd ..

go fmt  .\...
go tool vet .
go test .\...

genny -in pipe.go -out ../pipe-m.gen -pkg pipe gen "Thing=Any"
@Echo.
@Echo compare ../pipe-m.* - 2 Differences @ top: generic.Any ...
@Echo.

@pause
