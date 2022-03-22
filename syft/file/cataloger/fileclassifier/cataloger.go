package fileclassifier

import (
	"github.com/anchore/syft/internal/log"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/source"
)

type Cataloger struct {
	classifiers []file.Classifier
}

func NewCataloger(classifiers []file.Classifier) (*Cataloger, error) {
	return &Cataloger{
		classifiers: classifiers,
	}, nil
}

func (i *Cataloger) Catalog(resolver source.FileResolver) (map[source.Coordinates][]file.Classification, error) {
	results := make(map[source.Coordinates][]file.Classification)

	numResults := 0
	for _, location := range source.AllRegularFiles(resolver) {
		for _, classifier := range i.classifiers {
			result, err := classifier.Classify(resolver, location)
			if err != nil {
				log.Warnf("file classification cataloger failed with class=%q at location=%+v: %+v", classifier.Class, location, err)
				continue
			}
			if result != nil {
				results[location.Coordinates] = append(results[location.Coordinates], *result)
				numResults++
			}
		}
	}
	log.Debugf("file classifier discovered %d results", numResults)

	return results, nil
}