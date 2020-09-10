/*
Copyright (c) SiteWhere, LLC. All rights reserved. http://www.sitewhere.com

The software in this package is published under the terms of the CPAL v1.0
license, a copy of which has been included with this distribution in the
LICENSE file.
*/

package action

import "github.com/sitewhere/swctl/pkg/install"

// CheckInstall is the action for check SiteWhere installation
type CheckInstall struct {
	cfg *Configuration

	// Use verbose mode
	Verbose bool
}

// NewCheckInstall constructs a new *Install
func NewCheckInstall(cfg *Configuration) *CheckInstall {
	return &CheckInstall{
		cfg:     cfg,
		Verbose: false,
	}
}

// Run executes the list command, returning a set of matches.
func (i *CheckInstall) Run() (*install.SiteWhereInstall, error) {
	if err := i.cfg.KubeClient.IsReachable(); err != nil {
		return nil, err
	}

	return &install.SiteWhereInstall{}, nil
}
