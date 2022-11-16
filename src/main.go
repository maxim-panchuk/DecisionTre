package main

import (
	dataservice "DecisionTre/src/dataService"
	decisiontree "DecisionTre/src/decisionTree"
	"fmt"
)

func main() {
	ar, sr := dataservice.Parse("../resources/data_train.txt", ";")

	dt := decisiontree.DecisionTree{SampleRange: sr, AttributeRange: ar}

	root := dt.Init(40)

	_, ast := dataservice.Parse("../resources/data_test.txt", ";")

	for _, sample := range ast.SampleList {
		actual := sample.Class
		pr, _ := dt.Predict(root, sample)

		fmt.Printf("actual: %v    prediction: %v   \n", actual, pr)

	}

}
