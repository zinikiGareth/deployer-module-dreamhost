package dns

import (
	"ziniki.org/deployer/coremod/pkg/corebottom"
	"ziniki.org/deployer/driver/pkg/driverbottom"
	"ziniki.org/deployer/driver/pkg/errorsink"
)

type CNAMEBlank struct {
}

func (c *CNAMEBlank) ShortDescription() string {
	return "dreamhost.CNAME[]"
}

func (c *CNAMEBlank) Find(tools *corebottom.Tools, loc *errorsink.Location, id corebottom.CoinId, named string) corebottom.FindCoin {
	return &cnameCreator{tools: tools, loc: loc, name: named, coin: id}
}

func (c *CNAMEBlank) Mint(tools *corebottom.Tools, loc *errorsink.Location, id corebottom.CoinId, named string, props map[driverbottom.Identifier]driverbottom.Expr, teardown corebottom.TearDown) corebottom.Ensurable {
	return &cnameCreator{tools: tools, loc: loc, name: named, coin: id, props: props}
}

var _ corebottom.Blank = &CNAMEBlank{}
