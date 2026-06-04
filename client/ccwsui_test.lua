local CCWSUI = require("ccwsui").new("testing")
local c = require("components")

--- @param ctx CCWSUI.Context
local function line(ctx)
	return c.Literal("Hello!!! You're click #" .. tostring(ctx.state.count or 0))
		:hexColor("#cba6f7")
		:wrapPadding(2, 2, 2, 2)
		:wrapTexture("plain"):pad():tintedHex("#1e1e2e")
		:wrapClickRegion(ctx, function(ev)
			local add = ev.shift and 10 or 1
			ctx.state.count = (ctx.state.count or 0) + add
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
