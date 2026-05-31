.PHONY: gen_icon
gen_icon:
	wails3 generate icons -input appicon.png -macfilename darwin/icons.icns -windowsfilename windows/icon.ico -iconcomposerinput appicon.icon -macassetdir darwin

.PHONY: gen_dev_icon
gen_dev_icon:
	wails3 generate icons -input appicon_dev.png -macfilename darwin/icons_dev.icns -windowsfilename windows/icon_dev.ico -iconcomposerinput appicon_dev.icon -macassetdir darwin
