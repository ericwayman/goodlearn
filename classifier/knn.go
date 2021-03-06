package classifier

import (
	"fmt"

	"github.com/amitkgupta/goodlearn/classifier/knnutilities"
	"github.com/amitkgupta/goodlearn/dataset/dataset"
	"github.com/amitkgupta/goodlearn/dataset/row"
	"github.com/amitkgupta/goodlearn/dataset/target"
)

func NewKNNClassifier(k int) (*kNNClassifier, error) {
	if k < 1 {
		return nil, newInvalidNumberOfNeighboursError(k)
	}

	return &kNNClassifier{k: k}, nil
}

type kNNClassifier struct {
	k            int
	trainingData *dataset.Dataset
}

func (classifier *kNNClassifier) Train(trainingData *dataset.Dataset) error {
	if !trainingData.AllFeaturesFloats {
		return newNonFloatFeaturesTrainingSetError()
	}

	if trainingData.NumRows() == 0 {
		return newEmptyTrainingDatasetError()
	}

	classifier.trainingData = trainingData
	return nil
}

func (classifier *kNNClassifier) Classify(testRow *row.Row) (target.Target, error) {
	trainingData := classifier.trainingData
	if trainingData == nil {
		return nil, newUntrainedClassifierError()
	}

	if testRow.NumFeatures != trainingData.NumFeatures {
		return nil, newRowLengthMismatchError(testRow.NumFeatures, trainingData.NumFeatures)
	}

	if !testRow.AllFeaturesFloats {
		return nil, newNonFloatFeaturesTestRowError()
	}

	nearestNeighbours, err := knnutilities.NewKNNTargetCollection(classifier.k)
	if err != nil {
		return nil, err
	}

	testRowFeatureValues := testRow.UnsafeFloatFeatureValues()

	for i := 0; i < trainingData.NumRows(); i++ {
		trainingRow, err := trainingData.Row(i)
		if err != nil {
			return nil, err
		}
		trainingRowFeatureValues := trainingRow.UnsafeFloatFeatureValues()

		distance := knnutilities.Euclidean(testRowFeatureValues, trainingRowFeatureValues, nearestNeighbours.MaxDistance())
		if distance < nearestNeighbours.MaxDistance() {
			nearestNeighbours.Insert(trainingRow.Target, distance)
		}
	}

	return nearestNeighbours.Vote()
}

type EmptyTrainingDatasetError struct{}

func (e EmptyTrainingDatasetError) Error() string {
	return "cannot train on an empty dataset"
}
func newEmptyTrainingDatasetError() EmptyTrainingDatasetError {
	return EmptyTrainingDatasetError{}
}

type UntrainedClassifierError struct{}

func (e UntrainedClassifierError) Error() string {
	return "cannot classify before training"
}
func newUntrainedClassifierError() UntrainedClassifierError {
	return UntrainedClassifierError{}
}

type InvalidNumberOfNeighboursError struct {
	k int
}

func (e InvalidNumberOfNeighboursError) Error() string {
	return fmt.Sprintf("invalid number of neighbours %d", e.k)
}
func newInvalidNumberOfNeighboursError(k int) InvalidNumberOfNeighboursError {
	return InvalidNumberOfNeighboursError{k}
}

type InvalidFloatFeatureDatasetError struct{}

func (e InvalidFloatFeatureDatasetError) Error() string {
	return "dataset invalid, has feature columns which aren't floats"
}
func newInvalidFloatFeatureDatasetError() InvalidFloatFeatureDatasetError {
	return InvalidFloatFeatureDatasetError{}
}

type RowLengthMismatchError struct {
	numTestRowFeatures     int
	numTrainingSetFeatures int
}

func (e RowLengthMismatchError) Error() string {
	return fmt.Sprintf("Test row has %d features, training set has %d", e.numTestRowFeatures, e.numTrainingSetFeatures)
}
func newRowLengthMismatchError(numTestRowFeatures, numTrainingSetFeatures int) RowLengthMismatchError {
	return RowLengthMismatchError{numTestRowFeatures, numTrainingSetFeatures}
}

type NonFloatFeaturesTrainingSetError struct {
	kNNClassifier
}

func (e NonFloatFeaturesTrainingSetError) Error() string {
	return "cannot train on dataset with some non-float features"
}
func newNonFloatFeaturesTrainingSetError() NonFloatFeaturesTrainingSetError {
	return NonFloatFeaturesTrainingSetError{}
}

type NonFloatFeaturesTestRowError struct{}

func (e NonFloatFeaturesTestRowError) Error() string {
	return "cannot classify row with some non-float features"
}
func newNonFloatFeaturesTestRowError() NonFloatFeaturesTestRowError {
	return NonFloatFeaturesTestRowError{}
}
