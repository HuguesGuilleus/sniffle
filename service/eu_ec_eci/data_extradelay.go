package eu_ec_eci

var extraDelay_2020_9226 = ExtraDelay{
	Code:  "C(2020)9226",
	CELEX: "32020D2200",
}
var extraDelay_2021_1121 = ExtraDelay{
	Code:  "C(2021)1121",
	CELEX: "32021D0360",
}
var extraDelay_2021_3879 = ExtraDelay{
	Code:  "C(2021)3879",
	CELEX: "32021D0944",
}

type extraDelayICE struct {
	ID         uint
	Code       string
	Name       string
	ExtraDelay []ExtraDelay
}

var extraDelayMap = func() map[uint][]ExtraDelay {
	m := make(map[uint][]ExtraDelay, len(extraDelayData))
	for _, extra := range extraDelayData {
		m[extra.ID] = extra.ExtraDelay
	}
	return m
}()

var extraDelayData = [...]extraDelayICE{
	{240, "2019/6", "The fast, fair and effective solution to climate change", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{51, "2019/7", "Cohesion policy for the equality of the regions...", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{242, "2019/9", "Ending the aviation fuel tax exemption in Europe", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{246, "2019/11", "A price for carbon to fight climate change", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{247, "2019/12", "Grow scientific progress: crops matter!", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{252, "2019/14", "Stop corruption in Europe at its root, ...", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{253, "2019/15", "Actions on Climate Emergency", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{254, "2019/16", "Save Bees and farmers! ...", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{255, "2020/1", "Stop Finning – Stop the trade", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121}},
	{269, "2020/2", "VOTERS WITHOUT BORDERS...", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{267, "2020/3", "Start Unconditional Basic Incomes (UBI) throughout the EU", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{298, "2020/4", "Freedom to share", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{368, "2020/5", "Right to Cure", []ExtraDelay{extraDelay_2020_9226, extraDelay_2021_1121, extraDelay_2021_3879}},
	{648, "2021/1", "Civil society initiative for a ban on biometric mass surveillance practices", []ExtraDelay{extraDelay_2021_1121, extraDelay_2021_3879}},
	{814, "2021/3", "Green Garden Roof Tops", []ExtraDelay{extraDelay_2021_3879}},
}
