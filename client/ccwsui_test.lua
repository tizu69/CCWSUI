local CCWSUI = require("ccwsui").new("testing", "ws://localhost:8080")
local c = require("components")
local cc = require("compoundcomponents")

--- @param ctx CCWSUI.Context
local function line(ctx)
	local user = ctx:s("user")
	local global = ctx:s("global")
	return c.Literal("Hello!!! You're click #" .. tostring(user.count or 0))
		:hexColor("#555555")
		:text("\nTimer " .. tostring(global.timer or 0))
		:wrapPadding(2, 2, 2, 2)
		:wrapTexture("plain"):pad()
		:wrapClickRegion(ctx, function(ev)
			local add = ev.shift and 10 or 1
			user.count = (user.count or 0) + add
		end)
end

function CCWSUI:render(ctx)
	CCWSUI:updateMetadata(ctx.client) {
		title = "Testing"
	}

	return c.StackV()
		:add(
			line(ctx)
			:wrapAlignCenter()
			:wrapMediaQuery()
			:minWidth(200)
			:wrapExpand()
		)
		:add(
			cc.CreateWindow("Hello!", "#888b79",
				c.StackV()
				:add(line(ctx))
				:add(line(ctx))
				:wrapPadding(8),
				cc.Tooltip("Reset click count",
					cc.CreateButton(ctx, "trash", function(ev)
						ctx:s("user").count = 0
					end)
				),
				c.ExpandH(c.Blank()),
				cc.Separator("h", 4),
				cc.Tooltip("Round down to nearest 10",
					cc.CreateButton(ctx, "chevrons", function(ev)
						ctx:s("user").count = math.floor(((ctx:s("user").count or 0) - 1) / 10) * 10
					end, 180)
				),
				cc.Tooltip("Round up to nearest 10",
					cc.CreateButton(ctx, "chevrons", function(ev)
						ctx:s("user").count = math.ceil(((ctx:s("user").count or 0) + 1) / 10) * 10
						if ctx:s("user").count == 0 then ctx:s("user").count = 0 end -- fix for -0
					end)
				)
			)
			:wrapConstrain(300, 0)
			:wrapAlignCenter()
		)
		:add(line(ctx)
			:wrapAlignCenter()
			:wrapExpand())
end

parallel.waitForAny(function() CCWSUI:run() end, function()
	while true do
		sleep(1)
		CCWSUI:getState("global").timer = (CCWSUI:getState("global").timer or 0) + 1
	end
end)
print(CCWSUI.error)
