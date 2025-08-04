package dns

import (
	"fmt"
	"log"

	"ziniki.org/deployer/coremod/pkg/corebottom"
	"ziniki.org/deployer/driver/pkg/driverbottom"
	"ziniki.org/deployer/driver/pkg/errorsink"
	"ziniki.org/deployer/driver/pkg/utils"
	"ziniki.org/deployer/modules/dreamhost/internal/env"
)

type cnameCreator struct {
	tools *corebottom.Tools

	loc   *errorsink.Location
	name  string
	coin  corebottom.CoinId
	props map[driverbottom.Identifier]driverbottom.Expr
}

func (cc *cnameCreator) Loc() *errorsink.Location {
	return cc.loc
}

func (cc *cnameCreator) ShortDescription() string {
	return "dreamhost.CNAME[" + cc.name + "]"
}

func (cc *cnameCreator) DumpTo(iw driverbottom.IndentWriter) {
	iw.Intro("dreamhost.CNAME")
	iw.AttrsWhere(cc)
	iw.TextAttr("named", cc.name)
	iw.EndAttrs()
}

func (cc *cnameCreator) CoinId() corebottom.CoinId {
	return cc.coin
}

func (cc *cnameCreator) DetermineInitialState(pres corebottom.ValuePresenter) {
	eq := cc.tools.Recall.ObtainDriver("dreamhost.DreamhostEnv")
	dhEnv, ok := eq.(*env.DreamhostEnv)
	if !ok {
		panic("could not cast env to DreamhostEnv")
	}
	pointsTo, err := dhEnv.FindRecord("CNAME", cc.name)
	if err != nil {
		panic(err)
	}
	// log.Printf("points to %s\n", pointsTo)

	if pointsTo != "" {
		pt, ok := utils.AsStringer(pointsTo)
		if !ok {
			panic("not a stringer")
		}
		model := &cnameModel{loc: cc.loc, name: cc.name, pointsTo: pt}
		pres.Present(model)
	} else {
		pres.NotFound()
	}
}

func (cc *cnameCreator) DetermineDesiredState(pres corebottom.ValuePresenter) {
	var pointsTo driverbottom.Expr
	seenErr := false
	for p, v := range cc.props {
		switch p.Id() {
		case "PointsTo":
			pointsTo = v
		default:
			cc.tools.Reporter.ReportAtf(p.Loc(), "invalid property for IAM policy: %s", p.Id())
			seenErr = true
		}
	}
	if !seenErr && pointsTo == nil {
		cc.tools.Reporter.ReportAtf(cc.loc, "no PointsTo property was specified for %s", cc.name)
	}

	pt, ok := cc.tools.Storage.EvalAsStringer(pointsTo)
	if !ok {
		panic("not a stringer")
	}

	model := &cnameModel{loc: cc.loc, name: cc.name, pointsTo: pt}
	pres.Present(model)
}

func (cc *cnameCreator) UpdateReality() {
	tmp := cc.tools.Storage.GetCoin(cc.coin, corebottom.DETERMINE_INITIAL_MODE)

	if tmp != nil {
		found := tmp.(*cnameModel)
		log.Printf("CNAME %s already exists\n", found.name)
		return
	}

	log.Printf("creating CNAME %s\n", cc.name)
	desired := cc.tools.Storage.GetCoin(cc.coin, corebottom.DETERMINE_DESIRED_MODE).(*cnameModel)

	created := &cnameModel{name: cc.name, loc: cc.loc, pointsTo: desired.pointsTo}

	eq := cc.tools.Recall.ObtainDriver("dreamhost.DreamhostEnv")
	dhEnv, ok := eq.(*env.DreamhostEnv)
	if !ok {
		panic("could not cast env to DreamhostEnv")
	}

	err := dhEnv.InsertDNSRecord(desired.name, "CNAME", desired.pointsTo.String())
	if err != nil {
		panic(err)
	}

	cc.tools.Storage.Bind(cc.coin, created)
}

func (cc *cnameCreator) TearDown() {
	tmp := cc.tools.Storage.GetCoin(cc.coin, corebottom.DETERMINE_INITIAL_MODE)

	if tmp == nil {
		log.Printf("CNAME %s already deleted\n", cc.name)
		return
	}

	found := tmp.(*cnameModel)
	log.Printf("need to remove a CNAME record for %s\n", cc.name)
	log.Printf("found DH CNAME = %v\n", found)

	eq := cc.tools.Recall.ObtainDriver("dreamhost.DreamhostEnv")
	dhEnv, ok := eq.(*env.DreamhostEnv)
	if !ok {
		panic("could not cast env to DreamhostEnv")
	}

	err := dhEnv.DeleteDNSRecord(found.name, "CNAME")
	if err != nil {
		panic(err)
	}
}

func (cc *cnameCreator) String() string {
	return fmt.Sprintf("EnsureCNAME[%s:%s]", "" /* eb.env.Region */, cc.name)
}

var _ corebottom.Ensurable = &cnameCreator{}
