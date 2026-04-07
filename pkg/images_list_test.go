package pkg

import (
	"testing"

	"github.com/byDimasik/helm-images/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImagesCollectImagesFromTemplateRecursesIntoListItems(t *testing.T) {
	imageClient := Images{
		Kind: []string{k8s.KindDeployment},
	}
	imageClient.SetLogger("info")

	manifest := `
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ConfigMap
    metadata:
      name: should-be-skipped
    data:
      test: data
  - apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: nested-deployment
    spec:
      template:
        spec:
          containers:
            - name: app
              image: example.com/app:1.0.0
`

	images, err := imageClient.collectImagesFromTemplate(manifest, nil)
	require.NoError(t, err)
	require.Len(t, images, 1)
	assert.Equal(t, "nested-deployment", images[0].Name)
	assert.Equal(t, []string{"example.com/app:1.0.0"}, images[0].Image)
}
