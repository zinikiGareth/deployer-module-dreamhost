package dns

import (
	"fmt"
	"log"

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

	log.Printf("dhEnv = %p\n", dhEnv)

	records, err := dhEnv.DNSRecordsFor(dnf.name)
	// dnf.domainsClient = awsEnv.Route53DomainsClient()
	// detail, err := dnf.domainsClient.GetDomainDetail(context.TODO(), &route53domains.GetDomainDetailInput{DomainName: &dnf.name})
	if err != nil {
		panic(err)
	}

	/*
		var hzid string
		for _, z := range zones.HostedZones {
			if *z.Name == dnf.name+"." {
				hzid = strings.Replace(*z.Id, "/hostedzone/", "", 1)
				log.Printf("found zone %s: %s\n", hzid, *z.Name)
			}
		}
		if hzid == "" {
			log.Fatalf("No hosted zone found for " + dnf.name)
		}
	*/
	model := CreateDomainModel(dnf.loc, records)
	pres.Present(model)
}

func (dnf *domainNameFinder) String() string {
	return fmt.Sprintf("FindDomainName[%s]", dnf.name)
}

var _ corebottom.FindCoin = &domainNameFinder{}
