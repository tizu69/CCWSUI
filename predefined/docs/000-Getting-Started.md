> [!CAUTION]
> CCWSUI is not an official Minecraft service. It is not approved by or
> associated with Mojang or Microsoft.

Welcome to the CCWSUI developer documentation!

CCWSUI is a UI framework that offers a simple, declarative way to build familiar
looking, interactive, and pixel-perfect web UIs, controlled entirely from within
your ComputerCraft code.

> [!NOTE]
> In theory, anything that can communicate as a client with a WebSocket can
> serve a CCWSUI, however, at the moment, only ComputerCraft has an official
> host implementation.

---

To get started with building a CCWSUI for your ComputerCraft scripts, first
generate a CCWSUI bundle from the [Downloads](page:010-Download) page.

The client is split into three files. The client itself, as ccwsui.lua, the
native components, as components.lua, and some compound components built in Lua
itself, as compoundcomponents.lua. To use the client, simply require it and
create a new instance, like so:

```lua
local CCWSUI = require("ccwsui").new("testing", nil)
local c = require("components")
local cc = require("compoundcomponents")
--[[br]]
CCWSUI:run()
```

The first argument is the slug you wish to claim. CCWSUI, being server-based,
relies on slugs, which serves as a unique identifier for your UI. The second
argument is the URL of the server you wish to connect to. Both may be nil,
in which case the default slug and address will be used.

`CCWSUI:run()` is a blocking function that will run the UI forever. You should
call this at the end of your script. If you wish to do other things apart from
rendering, use `parallel.waitForAny()` to run something else in parallel. If you
now run your script, you will see that a URL is printed to the console. This is
the URL you can use to access your UI. It may be `testing` as you specified, or
something else if the slug is already taken. Visit the URL in a browser to see
your UI.

Now, obviously, you have not actually implemented a UI yet. That is the purpose
of the `render` hook. This hook is called every time the UI needs to be rendered
for a client. For example, let's render a literal "Hello, world!" to the screen.

> [!NOTE]
> You must define the `render` hook before calling `CCWSUI:run()`. The `run`
> call would block forever, so your code setting the hook would never be called.

```lua
function CCWSUI:render(ctx)
	return c.Literal("Hello, world!")
end
```

And you're done, your first UI is complete! Of course, this is not very useful,
so let's get building!

[Next: Using Components](page:020-Using-Components)
