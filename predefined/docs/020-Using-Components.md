CCWSUI provides [a set of primitive components](page:030-Component-Gallery),
which are in components.lua. You can create compound components, made up of
these primitives. CCWSUI also ships with some built-in compound components in
the file compoundcomponents.lua.

If you've used Flutter before, the way CCWSUI components get used by you, as the
developer, may feel familiar to you.

Native components are created by calling a function on the `c` table. They will
return a `CCWSUI.Component`, which may have some methods for modifying it. The
render hook expects a return value of a `CCWSUI.Component`. Here is an example:

```lua
local c = require("components")
--[[br]]
function CCWSUI:render(ctx)
	return c.StackV()
		:add(
			c.Literal("Hello, world!")
		)
		:add(
			c.ExpandV(
				c.AlignCenter(
					c.Padding(10,
						c.Literal("Lorem ipsum")
						:hexColor("#555555")
					)
				)
			)
		)
end
```

The nesting that this component tree creates is a bit ugly, so you can "wrap"
anything in a native component, as long as that native component takes exactly
one child in its constructor. This means, the above can be written as:

```lua
local c = require("components")
--[[br]]
function CCWSUI:render(ctx)
	return c.StackV()
		:add(
			c.Literal("Hello, world!")
		)
		:add(
			c.Literal("Lorem ipsum")
			:hexColor("#555555")
			:wrapPadding(10)
			:wrapAlignCenter()
			:wrapExpandV()
		)
end
```

Because native components are just values, you can can create helper functions
(also referred to as compound components), which you can call from any other
component-rendering context:

```lua
local c = require("components")
--[[br]]
---@param ctx CCWSUI.Context
---@param props { color: string }
local function lorem(ctx, props)
	return c.Literal("Lorem ipsum")
		:hexColor(props.color)
		:wrapPadding(10)
		:wrapAlignCenter()
		:wrapExpandV()
end
--[[br]]
function CCWSUI:render(ctx)
	return c.StackV()
		:add(c.Literal("Hello, world!"))
		:add(lorem(ctx, { color = "#555555" }))
end
```

[Next: Using State](page:040-Using-State)
