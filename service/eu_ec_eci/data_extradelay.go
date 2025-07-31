package eu_ec_eci

import "github.com/HuguesGuilleus/sniffle/front/component"

var extraDelay_2020_9226 = component.Legal{
	Num:   "C(2020)9226",
	CELEX: "32020D2200",
}
var extraDelay_2021_1121 = component.Legal{
	Num:   "C(2021)1121",
	CELEX: "32021D0360",
}
var extraDelay_2021_3879 = component.Legal{
	Num:   "C(2021)3879",
	CELEX: "32021D0944",
}

type extraDelayICE struct {
	ID         uint
	Code       string
	Name       string
	ExtraDelay []component.Legal
}

var extraDelayMap = func() map[uint][]component.Legal {
	m := make(map[uint][]component.Legal, len(extraDelayData))
	for _, extra := range extraDelayData {
		m[extra.ID] = extra.ExtraDelay
	}
	return m
}()

var extraDelayData = [...]extraDelayICE{
	{240, "2019/6", "The fast, fair and effective solution to climate change", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{51, "2019/7", "Cohesion policy for the equality of the regions...", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{242, "2019/9", "Ending the aviation fuel tax exemption in Europe", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{246, "2019/11", "A price for carbon to fight climate change", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{247, "2019/12", "Grow scientific progress: crops matter!", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{252, "2019/14", "Stop corruption in Europe at its root, ...", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{253, "2019/15", "Actions on Climate Emergency", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{254, "2019/16", "Save Bees and farmers! ...", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{255, "2020/1", "Stop Finning – Stop the trade", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121}},
	{269, "2020/2", "VOTERS WITHOUT BORDERS...", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{267, "2020/3", "Start Unconditional Basic Incomes (UBI) throughout the EU", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{298, "2020/4", "Freedom to share", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{368, "2020/5", "Right to Cure", []component.Legal{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{648, "2021/1", "Civil society initiative for a ban on biometric mass surveillance practices", []component.Legal{extraDelay_2021_1121, extraDelay_2021_3879}},
	{814, "2021/3", "Green Garden Roof Tops", []component.Legal{extraDelay_2021_3879}},
}
