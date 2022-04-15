package helm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderTemplateE(t *testing.T) {
	res, err := RenderTemplateE(t, nil, "../../examples/helm-dependency-example", "", nil)
	require.NoError(t, err)
	require.Empty(t, res)
}
