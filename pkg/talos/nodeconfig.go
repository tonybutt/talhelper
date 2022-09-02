package talos

import (
	"github.com/budimanjojo/talhelper/pkg/config"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/generate"
	"github.com/talos-systems/talos/pkg/machinery/config/types/v1alpha1/machine"
)

func GenerateNodeConfigBytes(node *config.Nodes, input *generate.Input) ([]byte, error) {
	cfg, err := generateNodeConfig(node, input)
	if err != nil {
		return nil, err
	}
	return cfg.EncodeBytes()
}

func generateNodeConfig(node *config.Nodes, input *generate.Input) (*v1alpha1.Config, error) {
	var c *v1alpha1.Config
	var err error

	nodeInput, err := patchNodeInput(node, input)
	if err != nil {
		return nil, err
	}

	switch node.ControlPlane {
	case true:
		c, err = generate.Config(machine.TypeControlPlane, nodeInput)
		if err != nil {
			return nil, err
		}
	case false:
		c, err = generate.Config(machine.TypeWorker, nodeInput)
		if err != nil {
			return nil, err
		}
	}

	cfg := applyNodeOverride(node, c)

	return cfg, nil
}

func applyNodeOverride(node *config.Nodes, cfg *v1alpha1.Config) (*v1alpha1.Config) {
	cfg.MachineConfig.MachineNetwork.NetworkHostname = node.Hostname

	if len(node.Nameservers) != 0 {
		cfg.MachineConfig.MachineNetwork.NameServers = node.Nameservers
	}

	if node.DisableSearchDomain {
		cfg.MachineConfig.MachineNetwork.NetworkDisableSearchDomain = &node.DisableSearchDomain
	}

	if len(node.NetworkInterfaces) != 0 {
		cfg.MachineConfig.MachineNetwork.NetworkInterfaces = node.NetworkInterfaces
	}

	return cfg
}

func patchNodeInput(node *config.Nodes, input *generate.Input) (*generate.Input, error) {
	nodeInput := input
	if node.InstallDisk != "" {
		nodeInput.InstallDisk = node.InstallDisk
	}

	return nodeInput, nil
}
