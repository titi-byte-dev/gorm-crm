package version

// Estas variáveis são injectadas em compile-time via -ldflags.
// Em desenvolvimento ficam com os valores default abaixo.
// Em CI/CD o Makefile e o Dockerfile passam os valores reais.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)
