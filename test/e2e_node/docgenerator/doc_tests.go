package docgenerator

import (
	"fmt"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"gopkg.in/yaml.v3"
)


type TestCaseType string

const (
	POD_CONDITION_TEST TestCaseType = "POD_CONDITION_TEST"
	CONTAINER_LIFECYCLE_TEST TestCaseType = "CONTAINER_LIFECYCLE_TEST"
)

var Lock sync.Mutex
var TestOutputs map[TestCaseType][]interface{}

func init() {
	 TestOutputs = make(map[TestCaseType][]interface{})
}

type DocGenerator interface {
	Generate(builder *strings.Builder) error
}

type PodConditionsTestOutput struct {
	Title string `json:"title"`
	TestDescription string `json:"testdescription"`
	PodSpec v1.Pod `json:"podSpec"`
	SuccessfulStates []State `json:"successfulstates"`
	FailedStates []State `json:"failedstates"`
	FailureReason string `json:"failurereason"`
}

type State struct {
	State v1.PodConditionType `json:"state"`
	TimeTaken time.Time `json:"timetaken"`
}

func TitleLevel(title string, level int) string {
	return strings.Repeat("#", level) + " " + title + "\n"
}

func CodeBlock(block interface{}) string {
	return ""
}

func Paragraph(paragraph *string) {

}

func (pc *PodConditionsTestOutput) Generate(builder *strings.Builder) error {
	builder.WriteString("\n" + TitleLevel(pc.Title, 2) + "\n\n")

	y, err := yaml.Marshal(pc.PodSpec)

    if err != nil {
        return err
    }

	builder.WriteString(TitleLevel(pc.TestDescription, 3) + "\n\n")

    content := fmt.Sprintf(
        "## PodSpec\n\n```yaml\n%s```\n\n",
        string(y),
    )

	builder.WriteString(content)

	var successStates []string
	for _, success := range pc.SuccessfulStates {
		successStates = append(successStates, string(success.State))
	}

	if len(successStates) > 0 {
		builder.WriteString("The POD successfully transitions through these states\n")
		builder.WriteString("**" + strings.Join(successStates, " => ") + "**\n\n")
	}

	var failedStates []string
	for _, failed := range pc.FailedStates {
		failedStates = append(failedStates, string(failed.State))
	}

	if len(failedStates) > 0 {
		builder.WriteString("The POD fails to transitions through these states\n")
		builder.WriteString("**" + strings.Join(failedStates, " => ") + "**\n\n")
	}

	if pc.FailureReason != "" {
		builder.WriteString(fmt.Sprintf("**Reason of failure**: %s\n", pc.FailureReason))
	}

	return nil

}

type ContainerLifeCycleTestOutput struct {


}

func (pc *ContainerLifeCycleTestOutput) Generate() {

}

func AddTestOutput(testType TestCaseType, testOutput interface{}) {
	Lock.Lock()
    defer Lock.Unlock()
	tests, ok := TestOutputs[testType]
	if ok {
		tests = append(tests, testOutput)
		TestOutputs[testType] = tests
		return
	}
	var newTests []interface{}
	newTests = append(newTests, testOutput)
	TestOutputs[testType] = newTests
	fmt.Println(TestOutputs)
}
