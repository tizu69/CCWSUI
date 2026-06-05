local CCWSUI = require("ccwsui").new("testing")
local c = require("components")

--- @param ctx CCWSUI.Context
local function line(ctx)
	local user = ctx:s("user")
	return c.Literal("Hello!!! You're click #" .. tostring(user.count or 0))
		:hexColor("#555555")
		:wrapPadding(2, 2, 2, 2)
		:wrapTexture("plain"):pad()
		:wrapClickRegion(ctx, function(ev)
			local add = ev.shift and 10 or 1
			user.count = (user.count or 0) + add
		end)
		:wrapAlignCenter()
		:wrapExpand()
end

function CCWSUI:render(ctx)
	return c.StackV()
		:add(line(ctx))
		:add(line(ctx))
		:add(line(ctx))
		:add(line(ctx))
		:add(line(ctx))
end

function CCWSUI:ready(url)
	print("Access UI: " .. url)
end

CCWSUI:run()
print(CCWSUI.error)
