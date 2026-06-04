---@alias CCWSUI.Alignment number
local AlignmentStart = 0
local AlignmentCenter = .5
local AlignmentEnd = 1

---@alias CCWSUI.Direction string
local DirectionH = "h"
local DirectionV = "v"
local DirectionHV = "hv"

---@alias CCWSUI.Flip string
local FlipNone = ""
local FlipX = "x"
local FlipY = "y"

---@alias CCWSUI.Rotation number
local RotationU = 0
local RotationR = 90
local RotationD = 180
local RotationL = 270

---@alias CCWSUI.StackDirection string
local StackDirectionH = "h"
local StackDirectionV = "v"

local emptyarr = textutils.empty_json_array

local alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
local function nanoid(len)
	local id = ""
	for _ = 1, len or 21 do
		local r = math.random(1, #alphabet)
		id = id .. alphabet:sub(r, r)
	end
	return id
end

local function hexToByte(hex) return tonumber(hex, 16) or 0 end

---@param hex string
---@return integer, integer, integer, integer
local function colorFromHex(hex)
	hex = hex:gsub("^#", "")
	local len = #hex
	if len == 3 then
		return hexToByte(hex:sub(1, 1) .. hex:sub(1, 1)),
			hexToByte(hex:sub(2, 2) .. hex:sub(2, 2)),
			hexToByte(hex:sub(3, 3) .. hex:sub(3, 3)),
			0xff
	elseif len == 6 then
		return hexToByte(hex:sub(1, 2)),
			hexToByte(hex:sub(3, 4)),
			hexToByte(hex:sub(5, 6)),
			0xff
	elseif len == 8 then
		return hexToByte(hex:sub(1, 2)),
			hexToByte(hex:sub(3, 4)),
			hexToByte(hex:sub(5, 6)),
			hexToByte(hex:sub(7, 8))
	end
	return 0, 0, 0, 0xff
end


---@class CCWSUI.Component
---@field protected kind string
---@field protected props table<string, any>
---@field protected children CCWSUI.Component[]
local component = {}
local function makecomp() return setmetatable({}, { __index = component }) end

local components = {}

---@class Align : CCWSUI.Component
local Align = makecomp()

---@param x CCWSUI.Alignment
---@param y CCWSUI.Alignment
---@param child CCWSUI.Component
---@return Align
function components.Align(x, y, child)
	return setmetatable({
		kind = "align",
		props = { X = x, Y = y },
		children = { child },
	}, { __index = Align })
end

---@param x CCWSUI.Alignment
---@param y CCWSUI.Alignment
---@return Align
function component:wrapAlign(x, y) return components.Align(x, y, self) end

---@param x CCWSUI.Alignment
---@param child CCWSUI.Component
---@return Align
function components.AlignX(x, child)
	return setmetatable({
		kind = "align",
		props = { X = x, Y = AlignmentCenter },
		children = { child },
	}, { __index = Align })
end

---@param x CCWSUI.Alignment
---@return Align
function component:wrapAlignX(x) return components.AlignX(x, self) end

---@param y CCWSUI.Alignment
---@param child CCWSUI.Component
---@return Align
function components.AlignY(y, child)
	return setmetatable({
		kind = "align",
		props = { X = AlignmentCenter, Y = y },
		children = { child },
	}, { __index = Align })
end

---@param y CCWSUI.Alignment
---@return Align
function component:wrapAlignY(y) return components.AlignY(y, self) end

---@param child CCWSUI.Component
---@return Align
function components.AlignCenter(child)
	return setmetatable({
		kind = "align",
		props = { X = AlignmentCenter, Y = AlignmentCenter },
		children = { child },
	}, { __index = Align })
end

---@return Align
function component:wrapAlignCenter() return components.AlignCenter(self) end

---@class Blank : CCWSUI.Component
local Blank = makecomp()

---@return Blank
function components.Blank()
	return setmetatable({
		kind = "blank",
		props = {},
		children = emptyarr,
	}, { __index = Blank })
end

---@param w number
---@param h number
---@return Blank
function components.Filler(w, h)
	return setmetatable({
		kind = "blank",
		props = { W = w, H = h },
		children = emptyarr,
	}, { __index = Blank })
end

---@class ClickRegion : CCWSUI.Component
local ClickRegion = makecomp()

---@class CCWSUI.ClickRegionEvent
---@field shift boolean
---@field ctrl boolean
---@field alt boolean

---@param ctx CCWSUI.Context
---@param event CCWSUI.Handler<CCWSUI.ClickRegionEvent>|string
---@param child CCWSUI.Component
---@return ClickRegion
function components.ClickRegion(ctx, event, child)
	local eventid = type(event) == "string" and event or nanoid()
	if type(event) == "function" then ctx:addHandler(eventid, event) end
	return setmetatable({
		kind = "clickregion",
		props = { Event = eventid },
		children = { child },
	}, { __index = ClickRegion })
end

---@param ctx CCWSUI.Context
---@param event CCWSUI.Handler<CCWSUI.ClickRegionEvent>|string
---@return ClickRegion
function component:wrapClickRegion(ctx, event) return components.ClickRegion(ctx, event, self) end

---@class Constrain : CCWSUI.Component
local Constrain = makecomp()

---@param w number
---@param h number
---@param child CCWSUI.Component
---@return Constrain
function components.Constrain(w, h, child)
	return setmetatable({
		kind = "constrain",
		props = { W = w, H = h },
		children = { child },
	}, { __index = Constrain })
end

---@param w number
---@param h number
---@return Constrain
function component:wrapConstrain(w, h) return components.Constrain(w, h, self) end

---@class Expanded : CCWSUI.Component
local Expanded = makecomp()

---@param child CCWSUI.Component
---@return Expanded
function components.Expand(child)
	return setmetatable({
		kind = "expanded",
		props = { Direction = DirectionHV },
		children = { child },
	}, { __index = Expanded })
end

---@return Expanded
function component:wrapExpand() return components.Expand(self) end

---@param child CCWSUI.Component
---@return Expanded
function components.ExpandH(child)
	return setmetatable({
		kind = "expanded",
		props = { Direction = DirectionH },
		children = { child },
	}, { __index = Expanded })
end

---@return Expanded
function component:wrapExpandH() return components.ExpandH(self) end

---@param child CCWSUI.Component
---@return Expanded
function components.ExpandV(child)
	return setmetatable({
		kind = "expanded",
		props = { Direction = DirectionV },
		children = { child },
	}, { __index = Expanded })
end

---@return Expanded
function component:wrapExpandV() return components.ExpandV(self) end

---@class FollowMouse : CCWSUI.Component
local FollowMouse = makecomp()

---@param x CCWSUI.Alignment
---@param y CCWSUI.Alignment
---@param child CCWSUI.Component
---@return FollowMouse
function components.FollowMouse(x, y, child)
	return setmetatable({
		kind = "followmouse",
		props = { X = x, Y = y },
		children = { child },
	}, { __index = FollowMouse })
end

---@param x CCWSUI.Alignment
---@param y CCWSUI.Alignment
---@return FollowMouse
function component:wrapFollowMouse(x, y) return components.FollowMouse(x, y, self) end

function FollowMouse:flipIfOverflow()
	self.props.FlipIfOverflow = true
	return self
end

---@class Icon : CCWSUI.Component
local Icon = makecomp()

---@param icon string
---@return Icon
function components.Icon(icon)
	return setmetatable({
		kind = "icon",
		props = { Icon = icon },
		children = emptyarr,
	}, { __index = Icon })
end

---@param icon string
---@return Icon
function component:wrapIcon(icon) return components.Icon(icon) end

---@param r number
---@param g number
---@param b number
---@param a number
function Icon:tinted(r, g, b, a)
	self.props.Icon = self.props.Icon .. string.format(";tint=#%02x%02x%02x%02x",
		r, g, b, a)
	return self
end

---@param hex string
function Icon:tintedHex(hex)
	return self:tinted(colorFromHex(hex))
end

--- @param rot CCWSUI.Rotation
function Icon:rotate(rot)
	self.props.Icon = self.props.Icon .. ";rotate=" .. tonumber(rot)
	return self
end

--- @param flip CCWSUI.Flip
function Icon:flip(flip)
	self.props.Icon = self.props.Icon .. ";flip=" .. flip
	return self
end

function Icon:shadow()
	self.props.Shadow = true
	return self
end

---@class ItemTexture : CCWSUI.Component
local ItemTexture = makecomp()

---@param item string
---@return ItemTexture
function components.ItemTexture(item)
	return setmetatable({
		kind = "itemtexture",
		props = { Item = item },
		children = emptyarr,
	}, { __index = ItemTexture })
end

---@class Literal : CCWSUI.Component
local Literal = makecomp()

---@param text string
---@return Literal
function components.Literal(text)
	return setmetatable({
		kind = "literal",
		props = { Pieces = { { Text = text, Color = { R = 255, G = 255, B = 255, A = 255 } } } },
		children = emptyarr,
	}, { __index = Literal })
end

---@param text string
function Literal:text(text)
	table.insert(self.props.Pieces, { Text = text, Color = { R = 255, G = 255, B = 255, A = 255 } })
	return self
end

---@param r number 0-255
---@param g number 0-255
---@param b number 0-255
---@param a number 0-255
function Literal:color(r, g, b, a)
	self.props.Pieces[#self.props.Pieces].Color = { R = r, G = g, B = b, A = a }
	return self
end

---@param hex string
function Literal:hexColor(hex)
	return self:color(colorFromHex(hex))
end

function Literal:shadow()
	self.props.Pieces[#self.props.Pieces].Shadow = true
	return self
end

function Literal:wrap()
	self.props.Wrap = true
	return self
end

---@param alignment CCWSUI.Alignment
function Literal:align(alignment)
	self.props.Alignment = alignment
	return self
end

---@class Overlay : CCWSUI.Component
local Overlay = makecomp()

---@return Overlay
function components.Overlay()
	return setmetatable({
		kind = "overlay",
		props = {},
		children = emptyarr,
	}, { __index = Overlay })
end

---@param child CCWSUI.Component
function Overlay:add(child)
	if self.children == emptyarr then self.children = {} end
	table.insert(self.children, child)
	return self
end

function Overlay:requireHover()
	self.props.HoverRequired = true
	return self
end

---@class Padding : CCWSUI.Component
local Padding = makecomp()

---@param t number
---@param l number
---@param b number
---@param r number
---@param child CCWSUI.Component
---@return Padding
function components.Padding(t, l, b, r, child)
	return setmetatable({
		kind = "padding",
		props = { T = t, L = l, B = b, R = r },
		children = { child },
	}, { __index = Padding })
end

---@param t number
---@param l number
---@param b number
---@param r number
---@return Padding
function component:wrapPadding(t, l, b, r) return components.Padding(t, l, b, r, self) end

---@class Scroll : CCWSUI.Component
local Scroll = makecomp()

---@param id string
---@param direction CCWSUI.Direction
---@param child CCWSUI.Component
---@return Scroll
function components.Scroll(id, direction, child)
	return setmetatable({
		kind = "scroll",
		props = { ID = id, Direction = direction, Step = 8 },
		children = { child },
	}, { __index = Scroll })
end

---@param id string
---@param direction CCWSUI.Direction
---@return Scroll
function component:wrapScroll(id, direction) return components.Scroll(id, direction, self) end

---@param step number
function Scroll:step(step)
	self.props.Step = step
	return self
end

---@class Stack : CCWSUI.Component
local Stack = makecomp()

---@return Stack
function components.StackH()
	return setmetatable({
		kind = "stack",
		props = { Direction = StackDirectionH },
		children = emptyarr,
	}, { __index = Stack })
end

---@return Stack
function components.StackV()
	return setmetatable({
		kind = "stack",
		props = { Direction = StackDirectionV },
		children = emptyarr,
	}, { __index = Stack })
end

---@param child CCWSUI.Component
function Stack:add(child)
	if self.children == emptyarr then self.children = {} end
	table.insert(self.children, child)
	return self
end

---@param gap number
function Stack:gap(gap)
	self.props.Padding = gap
	return self
end

---@class Texture : CCWSUI.Component
local Texture = makecomp()

---@param tex string
---@param child CCWSUI.Component
---@return Texture
function components.Texture(tex, child)
	return setmetatable({
		kind = "texture",
		props = { Tex = tex },
		children = { child },
	}, { __index = Texture })
end

---@param tex string
---@return Texture
function component:wrapTexture(tex) return components.Texture(tex, self) end

---@param r number
---@param g number
---@param b number
---@param a number
function Texture:tinted(r, g, b, a)
	self.props.Tex = self.props.Tex .. string.format(";tint=#%02x%02x%02x%02x",
		r, g, b, a)
	return self
end

---@param hex string
function Texture:tintedHex(hex)
	return self:tinted(colorFromHex(hex))
end

---@param sr number
---@param sg number
---@param sb number
---@param sa number
---@param dr number
---@param dg number
---@param db number
---@param da number
function Texture:remap(sr, sg, sb, sa, dr, dg, db, da)
	self.props.Tex = self.props.Tex ..
		string.format(";#%02x%02x%02x%02x=#%02x%02x%02x%02x",
			sr, sg, sb, sa, dr, dg, db, da)
	return self
end

---@param src string
---@param dst string
function Texture:remapHex(src, dst)
	local sr, sg, sb, sa = colorFromHex(src)
	local dr, dg, db, da = colorFromHex(dst)
	return self:remap(sr, sg, sb, sa, dr, dg, db, da)
end

function Texture:pad()
	self.props.Pad = true
	return self
end

return components
