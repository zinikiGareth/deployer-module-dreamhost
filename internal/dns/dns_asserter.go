package dns

import (
	"ziniki.org/deployer/driver/pkg/driverbottom"
	"ziniki.org/deployer/modules/dreamhost/internal/env"
)

func AssertDNSRecord(tools *driverbottom.CoreTools, zone, key, value string) error {
	dhenv := tools.Recall.ObtainDriver("dreamhost.DreamhostEnv").(*env.DreamhostEnv)
	return dhenv.AssertDNSRecord(zone, key, value)
}
