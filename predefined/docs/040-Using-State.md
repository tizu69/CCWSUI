CCWSUI comes with a built-in reactive state system - if state changes, all the
clients reading from this state will automatically rerender their UI. Due to
CCWSUI's nature of being server-sided, state is managed using state keys. This
means multiple clients can read from multiple state tables, possibly the same
ones.

State can be accessed either from the CCWSUI instance, or the context given to
CCWSUI:render(ctx). As such, the following two examples modify the same state:

```lua
function CCWSUI:render(ctx)
    local state = ctx:s("myState")
    return c.Literal(state.foo)
end
```

```lua
parallel.waitForAny(
	function() CCWSUI:run() end,
	function()
		local state = CCWSUI:getState("myState")
		print(state.foo)
	end
)
```

> [!NOTE]
> You can not modify state from within a render function, ever. Not even using
> CCWSUI:getState. This is because modifying state will cause a rerender, which
> will call the render function again, which will modify state again. This would
> cause an infinite loop. It is safe to modify state from within a callback
> function however, such as a click handler.

---

ctx:s has two special handlers for nil keys, and the literal key "user", for the
current client and user state. Client IDs are ephemeral, and will not be the
same between reloads of the browser. User IDs are persistent, and will be the
same between reloads. Moreover, a user may have multiple clients open at a time,
and thus multiple client IDs.

```lua
ctx:s() == CCWSUI:getState("client:" .. ctx.client)
ctx:s("user") == CCWSUI:getState("user:" .. ctx.user)
```

> [!CAUTION]
> Never ever expose the user ID to the client, nor allow the user ID alone to
> be the deciding factor for things like authorization. The user ID is not a
> secret, and can be easily spoofed by a malicious client.

---

If a client reads from a state key using ctx:s, it will be marked as dependent
on that state key. If the state key is modified, all clients that are currently
dependent will be rerendered. If, during another rerender, it no longer uses the
state key, it will be unmarked as dependent and will no longer be rerendered.

> [!NOTE]
> An assignment will only cause a rerender if the direct root table of the key
> is modified. For example, `ctx:s("room:1").messages = {}` will cause a
> rerender, but `ctx:s("room:1").messages[1] = "Hello"` will not.

This behaviour allows you to easily build things like chat rooms. Each chat gets
a state key, the clients that read from that chat use the key to list all the
messages, and when a new message is added, all clients that are currently
reading from that chat will be rerendered.

```lua
function CCWSUI:render(ctx)
	local room = ctx:s().room -- room the client is in
	local msgs = ctx:s("room:" .. room).messages
	local msg = msgs[#msgs] or "None, yet!"
	return c.Literal("Last message: " .. msg)
end
--[[br]]
-- in some other place...
local room = CCWSUI:getState("room:1")
table.insert(room.messages, "Hello, world!")
room.messages = room.messages -- trigger rerender
```
