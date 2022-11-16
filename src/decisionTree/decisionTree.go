package decisiontree

import (
	dataservice "DecisionTre/src/dataService"
	"math"
	"math/rand"
	"strconv"
)

type Node struct {
	attributeIdx int
	children     map[string]*Node
	class        string
}

type DecisionTree struct {
	SampleRange    dataservice.SampleRange
	AttributeRange dataservice.AttributeRange
	root           *Node
}

func newNode() *Node {
	node := Node{children: make(map[string]*Node)}
	return &node
}

func newLeaf(class string) *Node {
	leaf := Node{class: class}
	return &leaf
}

var selectedAttributes map[int][]string
var maxDepth int

func (dt DecisionTree) Init(md int) *Node {
	selectedAttributes = dt.selectAttributes()
	maxDepth = md
	dt.root = newNode()
	dt.root.children = build(dt.SampleRange, dt.root, 1)
	return dt.root
}

func (dt DecisionTree) selectAttributes() map[int][]string {
	length := len(dt.AttributeRange.AttributeMap)
	amount := int(math.Sqrt(float64(length)))

	attrs := make(map[int][]string, amount)

	for i := 0; i < amount; i++ {
		key := rand.Intn(length)
		attrs[key] = dt.AttributeRange.AttributeMap[key]
	}

	return attrs
}

func build(samples dataservice.SampleRange, predecessor *Node, currDepth int) map[string]*Node {
	if predecessor == nil {
		return nil
	}

	children := make(map[string]*Node)
	attrIdx := findBestAttr(samples)

	predecessor.attributeIdx = attrIdx

	attributeVals := samples.AttrsById(attrIdx)
	uniqAttrVals := dataservice.UniqueList(attributeVals)

	for _, attrVal := range uniqAttrVals {
		selectedSamples := samples.GetWhereEq(attrIdx, attrVal)
		selectedClasses := selectedSamples.GetClasses()
		selectedClassesUniq := dataservice.UniqueList(selectedClasses)

		if len(selectedClassesUniq) == 1 {
			children[attrVal] = newLeaf(selectedClassesUniq[0])
		} else if currDepth > maxDepth {
			maxClasses := 0
			majorClass := ""

			for _, class := range selectedClassesUniq {
				classFreq := frequency(selectedClasses, class)
				if classFreq > maxClasses {
					maxClasses = classFreq
					majorClass = class
				}
			}
			children[attrVal] = newLeaf(majorClass)
		} else {
			newNode := newNode()
			newNode.children = build(selectedSamples, newNode, currDepth+1)
			children[attrVal] = newNode
		}
	}

	return children
}

func (dt DecisionTree) Predict(root *Node, sample dataservice.Sample) (class string, err error) {
	current := root
	for current.children != nil {
		attributeIdx := current.attributeIdx
		attrVal := sample.AttrValMap[attributeIdx]
		if _, ok := current.children[attrVal]; ok {
			current = current.children[attrVal]
		} else {
			min := 2147483647.0
			var foundKey string
			for key := range current.children {
				kInt, _ := strconv.Atoi(key)
				valInt, _ := strconv.Atoi(attrVal)
				diff := math.Abs(float64(kInt - valInt))

				if diff < min {
					foundKey = key
					min = diff
				}
			}
			current = current.children[foundKey]
		}
	}
	return current.class, nil
}

func findBestAttr(samples dataservice.SampleRange) int {
	classes := samples.GetClasses()
	uniqueClasses := dataservice.UniqueList(classes)

	infoT := infoT(classes, uniqueClasses)

	attrIdx := 0
	maxGainRatio := 0.0

	for key := range selectedAttributes {
		infoX := infoX(samples, key)
		splitInfoX := splitInfoX(samples, key)
		gainRatio := (infoT - infoX) / splitInfoX

		if gainRatio > maxGainRatio {
			maxGainRatio = gainRatio
			attrIdx = key
		}
	}

	return attrIdx
}

func infoT(classes []string, uniqueClasses []string) (result float64) {
	for _, class := range uniqueClasses {
		freq := frequency(classes, class)
		div := float64(freq) / float64(len(classes))
		result -= div * (math.Log2(div))
	}

	return
}

func infoX(samples dataservice.SampleRange, attributeKey int) (result float64) {
	attributeVals := samples.AttrsById(attributeKey)
	uniqAttrVals := dataservice.UniqueList(attributeVals)

	for _, attrVal := range uniqAttrVals {
		selectedSamples := samples.GetWhereEq(attributeKey, attrVal)
		selectedClasses := selectedSamples.GetClasses()
		selectedClassesUniq := dataservice.UniqueList(selectedClasses)

		result += float64(len(selectedSamples.SampleList)) / float64(len(samples.SampleList)) * infoT(selectedClasses, selectedClassesUniq)
	}

	return
}

func splitInfoX(samples dataservice.SampleRange, attributeKey int) (result float64) {
	attributeVals := samples.AttrsById(attributeKey)
	uniqAttrVals := dataservice.UniqueList(attributeVals)

	for _, attrVal := range uniqAttrVals {
		selectedSamples := samples.GetWhereEq(attributeKey, attrVal)
		div := float64(len(selectedSamples.SampleList)) / float64(len(samples.SampleList))
		result -= div * math.Log2(div)
	}

	return
}

func frequency(values []string, target string) (count int) {
	for _, val := range values {
		if val == target {
			count++
		}
	}

	return
}
