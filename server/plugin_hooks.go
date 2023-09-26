package main

import "fmt"

func (p *Plugin) OnActivate() error {
	return nil
}

func (p *Plugin) OnDeactivate() error {
	return nil
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return fmt.Errorf("failed to load plugin configuration: %w", err)
	}

	p.setConfiguration(configuration)

	if configuration.Backend == BackendPDFTron {
		p.backend = NewPDFTronBackend()
	}
	if configuration.Backend == BackendGotenberg {
		p.backend = NewGotenbergBackend()
	}

	if p.backend == nil {
		return fmt.Errorf("backend not selected")
	}

	if err := p.backend.Init(*configuration); err != nil {
		return fmt.Errorf("error initializing backend: %w", err)
	}

	return nil
}
