package main

import (
	"ziniki.org/deployer/driver/pkg/driverbottom"
	"ziniki.org/deployer/modules/dreamhost/pkg/dhmod"
)

func ProvideTestRunner(runner driverbottom.TestRunner) error {
	return dhmod.ProvideTestRunner(runner)
}

func RegisterWithDriver(deployer driverbottom.Driver) error {
	return dhmod.RegisterWithDriver(deployer)
}
