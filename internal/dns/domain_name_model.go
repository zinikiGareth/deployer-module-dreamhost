package dns

import (
	"ziniki.org/deployer/driver/pkg/driverbottom"
	"ziniki.org/deployer/driver/pkg/errorsink"
	"ziniki.org/deployer/modules/dreamhost/internal/env"
)

type domainModel struct {
	loc     *errorsink.Location
	records []*env.DNSRecord
	// hzid string
}

func (d *domainModel) Loc() *errorsink.Location {
	return d.loc
}

func (d *domainModel) ShortDescription() string {
	return "DomainName[" /* + d.hzid*/ + "]"
}

func (d *domainModel) DumpTo(to driverbottom.IndentWriter) {
	to.Intro("DomainName")
	to.AttrsWhere(d)
	// to.TextAttr("hzid", d.hzid)
	to.EndAttrs()
}

// func (d *domainModel) HostedZoneId() string {
// 	return d.hzid
// }

func CreateDomainModel(loc *errorsink.Location, records []*env.DNSRecord) *domainModel {
	return &domainModel{loc: loc, records: records}
}

func (dnf *domainModel) ObtainMethod(name string) driverbottom.Method {
	switch name {
	/*
		case "zoneId":
			return &zoneIdMethod{}
	*/
	}
	return nil
}

/*
type zoneIdMethod struct {
}

func (a *zoneIdMethod) Invoke(s driverbottom.RuntimeStorage, on driverbottom.Expr, args []driverbottom.Expr) any {
	e := on.Eval(s)
	model, ok := e.(*domainModel)
	if !ok {
		panic(fmt.Sprintf("zoneId can only be called on a domain, not a %T", e))
	}
	if len(args) != 0 {
		panic("invalid number of arguments")
	}
	return model.hzid
}
*/

var _ driverbottom.Describable = &domainModel{}

var _ driverbottom.HasMethods = &domainModel{}

// var _ ExportedDomain = &domainModel{}
