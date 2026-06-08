local c = require("components")

local components = {}

--- @param text string|Literal
--- @param child CCWSUI.Component
--- @return CCWSUI.Component
function components.Tooltip(text, child)
	text = type(text) == "string" and c.Literal(text) or text --[[@as Literal]]
	return c.Overlay()
		:add(child)
		:add(
			c.FollowMouse(0, 1,
				text
				:wrapPadding(5, 3)
				:wrapTexture("tooltip")
				:wrapPadding(4, 0)
			):flipIfOverflow()
		):requireHover()
end

--- @param title string Title text
--- @param color string Title background color
--- @param child CCWSUI.Component
--- @param ... CCWSUI.Component The buttons to add to the bottom of the window
--- @return CCWSUI.Component
function components.CreateWindow(title, color, child, ...)
	return c.StackV()
		:add(
			c.Literal(title)
			:wrapAlignCenter()
			:wrapPadding(10, 2)
			:wrapTexture("title"):tintedHex(color)
		)
		:add(
			child
			:wrapTexture("content-inset")
			:wrapPadding(12, -1)
			:wrapTexture("checkerboard")
			:wrapTexture("sided-inset"):pad()
			:wrapPadding(1)
			:wrapTexture("#000000")
			:wrapPadding(2, -1)
		)
		:add(
			c.StackH():add(...):gap(4):align(1)
			:wrapPadding(7, 6)
			:wrapTexture("title")
		)
end

--- @param ctx CCWSUI.Context
--- @param icon string
--- @param callback string|CCWSUI.Handler<CCWSUI.ClickRegionEvent>
--- @param rotate CCWSUI.Rotation|nil
--- @return CCWSUI.Component
function components.CreateButton(ctx, icon, callback, rotate)
	rotate = rotate or 0
	return c.Icon(icon):rotate(rotate)
		:wrapTexture("plain-outset"):pad()
		:wrapClickRegion(ctx, callback)
end

--- @param direction CCWSUI.StackDirection
--- @param unpad number|nil The amount of padding to remove to fill parent
--- @return CCWSUI.Component
function components.Separator(direction, unpad)
	unpad = unpad or 0
	if direction == "h" then
		return c.Filler(2, 0)
			:wrapTexture("title-inset")
			:wrapPadding(0, -unpad)
	else
		return c.Filler(0, 2)
			:wrapTexture("title-inset")
			:wrapPadding(-unpad, 0)
	end
end

return components
