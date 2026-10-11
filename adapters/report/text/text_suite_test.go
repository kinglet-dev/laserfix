package text_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kinglet-dev/core/domain/geom"
	"github.com/kinglet-dev/laserfix/app"
)

func TestText(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Text Suite")
}

// copyAt is a removed 40 mm copy at (x, y).
func copyAt(x, y float64) app.Change {
	return app.Change{Kind: app.Copy, Element: 2, Of: 1, At: geom.Point{X: x, Y: y}, Length: 40}
}
