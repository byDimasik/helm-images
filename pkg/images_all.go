package pkg

import (
	"fmt"
	"strings"

	"github.com/byDimasik/helm-images/pkg/errors"
	"github.com/byDimasik/helm-images/pkg/k8s"
)

type skipReleaseInfo struct {
	name      string
	namespace string
}

func (image *Images) GetAllImages() error {
	releases, err := image.getChartsFromReleases()
	if err != nil {
		return err
	}

	releases = releasesToSkip(image.releasesToSkip).filterRelease(releases)

	imagesFromAllRelease := make([]k8s.Images, 0)

	for _, release := range releases {
		image.log.Debugf("fetching the images from release '%s' of namespace '%s'", release.Name, release.Namespace)

		images := make([]*k8s.Image, 0)
		kubeKindTemplates := image.GetTemplates([]byte(release.Manifest))
		skips := image.GetResourcesToSkip()

		for _, kubeKindTemplate := range kubeKindTemplates {
			imagesFound, err := image.collectImagesFromTemplate(kubeKindTemplate, skips)
			if err != nil {
				return err
			}

			images = append(images, imagesFound...)
		}

		if len(images) == 0 {
			image.log.Infof("the release '%s' of namespace '%s' does not have any images", release.Name, release.Namespace)

			continue
		}

		output := image.setOutput(images)

		imagesFromAllRelease = append(imagesFromAllRelease, k8s.Images{ImagesFromRelease: output, NameSpace: release.Namespace})
	}

	return image.renderer.Render(imagesFromAllRelease)
}

func (image *Images) SetReleasesToSkips() error {
	const resourceLength = 2

	releasesToBeSkipped := make([]skipReleaseInfo, len(image.SkipReleases))

	for index, skipRelease := range image.SkipReleases {
		parsedRelease := strings.SplitN(skipRelease, "=", resourceLength)
		if len(parsedRelease) != resourceLength {
			return &errors.ImageError{Message: fmt.Sprintf("unable to parse release skip '%s'", skipRelease)}
		}

		releasesToBeSkipped[index] = skipReleaseInfo{name: parsedRelease[0], namespace: parsedRelease[1]}
	}

	image.releasesToSkip = releasesToBeSkipped

	return nil
}
