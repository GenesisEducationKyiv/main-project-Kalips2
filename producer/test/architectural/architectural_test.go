package architectural

import (
	"github.com/matthewmcnew/archtest"
	"testing"
)

func TestDomainLayer(t *testing.T) {
	t.Run("Domain layer should not depend on Application, Presenter, Infrastructure layers.", func(t *testing.T) {
		layers := append(applicationLayer, presenterLayer...)
		layers = append(layers, infrastructureLayer...)
		archtest.Package(t, domainLayer...).ShouldNotDependOn(layers...)
	})
}

func TestApplicationLayer(t *testing.T) {
	t.Run("Application layer should not depend on Presenter, Infrastructure layers.", func(t *testing.T) {
		layers := append(infrastructureLayer, presenterLayer...)
		archtest.Package(t, application).ShouldNotDependOn(layers...)
	})
}

func TestPresenterLayer(t *testing.T) {
	t.Run("Presenter layer should not depend on Infrastructure layer.", func(t *testing.T) {
		archtest.Package(t, presenterLayer...).ShouldNotDependOn(infrastructureLayer...)
	})
}

func TestInfrastructureLayer(t *testing.T) {
	t.Run("Infrastructure layer should not depend on Presenter layer.", func(t *testing.T) {
		archtest.Package(t, infrastructureLayer...).ShouldNotDependOn(presenterLayer...)
	})
}
