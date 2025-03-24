package main

import (
	"github.com/downace/adguardvpn-desktop/internal/adguard"
)

var allExclusionModes = []struct {
	Value  adguard.ExclusionMode
	TSName string
}{
	{adguard.ExclusionModeGeneral, "GENERAL"},
	{adguard.ExclusionModeSelective, "SELECTIVE"},
}
