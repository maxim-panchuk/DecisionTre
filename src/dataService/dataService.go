package dataservice

import (
	"log"
	"os"
	"strings"
)

type Sample struct {
	Id         int
	AttrValMap map[int]string
	Class      string
}

type SampleRange struct {
	SampleList []Sample
}

// type Attribute struct {
// 	Id     string
// 	Values []string
// }

type AttributeRange struct {
	AttributeMap map[int][]string
}

// Определим, какие значения аттрибутов встречаются в исходной выборке экземпляров
func (ar AttributeRange) defineAttrValues(sr *SampleRange) {
	for key := range ar.AttributeMap {
		allAttrVals := sr.AttrsById(key)
		uniqueList := UniqueList(allAttrVals)
		ar.AttributeMap[key] = uniqueList
	}
}

// Функция которая возвращает слайс всех значений конкретного аттрибута, встречающихся во всех экзмеплярах
func (sr SampleRange) AttrsById(idx int) []string {
	var result []string = make([]string, 0)

	for _, sample := range sr.SampleList {
		result = append(result, sample.AttrValMap[idx])
	}
	return result
}

func (sr SampleRange) GetClasses() []string {
	classes := make([]string, 0)
	for _, class := range sr.SampleList {
		classes = append(classes, class.Class)
	}
	return classes
}

func (sr SampleRange) GetWhereEq(attributeIdx int, attributeVal string) (sampleRange SampleRange) {
	sampleRange = SampleRange{}

	for _, sample := range sr.SampleList {
		if sample.AttrValMap[attributeIdx] == attributeVal {
			sampleRange.SampleList = append(sampleRange.SampleList, sample)
		}
	}
	return
}

func Parse(filename, separator string) (ar AttributeRange, sr SampleRange) {
	rawData, err := os.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	data := string(rawData)
	lines := strings.Split(data, "\n")

	attributeRange := parseAttributes(lines[0], separator)
	sampleRange := parseSampleRange(lines[1:], separator)

	attributeRange.defineAttrValues(sampleRange)

	return *attributeRange, *sampleRange
}

func parseAttributes(line, separator string) *AttributeRange {
	columns := strings.Split(line, separator)
	size := len(columns)

	var attributeRange AttributeRange = AttributeRange{make(map[int][]string)}

	for i := 1; i < size-1; i++ {
		attributeRange.AttributeMap[i] = make([]string, 0)
	}

	return &attributeRange
}

func parseSampleRange(lines []string, separator string) *SampleRange {

	var sampleRange SampleRange = SampleRange{make([]Sample, 0)}

	for i, line := range lines {
		sample := parseSample(line, separator, i+1)
		sampleRange.SampleList = append(sampleRange.SampleList, *sample)
	}

	return &sampleRange
}

func parseSample(line string, separator string, idx int) *Sample {
	columns := strings.Split(line, separator)
	size := len(columns)

	var sample Sample = Sample{Id: idx}
	m := make(map[int]string)

	for i := 1; i < size-1; i++ {
		m[i] = columns[i]
	}

	sample.AttrValMap = m
	sample.Class = string(columns[size-1])[:1]

	return &sample
}

// Функция вернет set из list
func UniqueList(list []string) []string {

	unique := make([]string, 0)

	for _, item := range list {
		if !contains(unique, item) {
			unique = append(unique, item)
		}
	}

	return unique
}

func contains(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
