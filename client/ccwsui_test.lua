local CCWSUI = require("ccwsui").new("testing")
local c = require("components")

local visits = 0
function CCWSUI:render(client)
	visits = visits + 1
	return c.Literal("Hello!!! You're visit #" .. tostring(visits))
		:hexColor("#cba6f7")
		:wrapClickRegion("foobar")
		:wrapTexture("")
		:wrapAlignCenter()
		:wrapExpand()
end

function CCWSUI:ready(url)
	print("Access UI: " .. url)
end

CCWSUI:run()
print(CCWSUI.error)
