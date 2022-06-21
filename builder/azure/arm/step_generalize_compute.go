package arm

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-azure/builder/azure/common/constants"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepGeneralizeCompute struct {
	client           *AzureClient
	generalizeVM     func(resourceGroupName, computeName string) error
	say              func(message string)
	shouldGeneralize func() bool
	error            func(e error)
}

func NewStepGeneralizeCompute(client *AzureClient, ui packersdk.Ui, config *Config) *StepGeneralizeCompute {
	var step = &StepGeneralizeCompute{
		client: client,
		say: func(message string) { ui.Say(message) },
		shouldGeneralize: func() bool { return config.Generalize },
		error:  func(e error) { ui.Error(e.Error()) },
	}

	step.generalizeVM = step.generalize
	return step
}

func (s *StepGeneralizeCompute) generalize(resourceGroupName string, computeName string) error {
	_, err := s.client.Generalize(context.TODO(), resourceGroupName, computeName)
	if err != nil {
		s.say(s.client.LastError.Error())
	}
	return err
}

func (s *StepGeneralizeCompute) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if !s.shouldGeneralize() {
		return multistep.ActionContinue
	}
	s.say("Generalizing machine ...")

	var resourceGroupName = state.Get(constants.ArmResourceGroupName).(string)
	var computeName = state.Get(constants.ArmComputeName).(string)

	s.say(fmt.Sprintf(" -> ResourceGroupName : '%s'", resourceGroupName))
	s.say(fmt.Sprintf(" -> ComputeName       : '%s'", computeName))

	err := s.generalizeVM(resourceGroupName, computeName)

	return processStepResult(err, s.error, state)
}

func (*StepGeneralizeCompute) Cleanup(multistep.StateBag) {
}
