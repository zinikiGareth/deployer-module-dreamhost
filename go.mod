module ziniki.org/deployer/modules/dreamhost

go 1.24.3

require (
	ziniki.org/deployer/coremod v0.0.0
	ziniki.org/deployer/driver v0.0.0
)

replace ziniki.org/deployer/driver => ../deployer/driver
replace ziniki.org/deployer/coremod => ../deployer/coremod
