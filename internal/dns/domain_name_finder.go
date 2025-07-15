package dns

import (
	"fmt"

	"ziniki.org/deployer/coremod/pkg/corebottom"
	"ziniki.org/deployer/driver/pkg/driverbottom"
	"ziniki.org/deployer/driver/pkg/errorsink"
	"ziniki.org/deployer/modules/dreamhost/internal/env"
)

type ExportedDomain interface {
	HostedZoneId() string
}

type domainNameFinder struct {
	tools *corebottom.Tools

	loc  *errorsink.Location
	name string
	coin corebottom.CoinId
}

func (dnf *domainNameFinder) Loc() *errorsink.Location {
	return dnf.loc
}

func (dnf *domainNameFinder) ShortDescription() string {
	return "dreamhost.DomainName[" + dnf.name + "]"
}

func (dnf *domainNameFinder) DumpTo(iw driverbottom.IndentWriter) {
	iw.Intro("dreamhost.DomainName[")
	iw.AttrsWhere(dnf)
	iw.TextAttr("named", dnf.name)
	iw.EndAttrs()
}

func (dnf *domainNameFinder) CoinId() corebottom.CoinId {
	return dnf.coin
}

func (dnf *domainNameFinder) DetermineInitialState(pres corebottom.ValuePresenter) {
	eq := dnf.tools.Recall.ObtainDriver("dreamhost.DreamhostEnv")
	dhEnv, ok := eq.(*env.DreamhostEnv)
	if !ok {
		panic("could not cast env to DreamhostEnv")
	}

	records, err := dhEnv.DNSRecordsFor(dnf.name)
	if err != nil {
		panic(err)
	}

	model := CreateDomainModel(dnf.loc, records)
	pres.Present(model)
}

func (dnf *domainNameFinder) String() string {
	return fmt.Sprintf("FindDomainName[%s]", dnf.name)
}

var _ corebottom.FindCoin = &domainNameFinder{}
