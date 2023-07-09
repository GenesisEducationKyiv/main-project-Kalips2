package architectural

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestDomainLayer(t *testing.T) {
	layers := append(applicationLayer, presenterLayer...)
	layers = append(layers, infrastructureLayer...)
	archtest.Package(t, domainLayer...).ShouldNotDependOn(layers...)
}

func TestApplicationLayer(t *testing.T) {
	layers := append(infrastructureLayer, presenterLayer...)
	archtest.Package(t, application).ShouldNotDependOn(layers...)
}

func TestPresenterLayer(t *testing.T) {
	archtest.Package(t, presenterLayer...).ShouldNotDependOn(infrastructureLayer...)
}

func TestInfrastructureLayer(t *testing.T) {
	archtest.Package(t, infrastructureLayer...).ShouldNotDependOn(presenterLayer...)
}
