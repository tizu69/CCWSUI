local expect = require "cc.expect"

--- @alias CCWSUI.Handler<T> fun(ev: T)

--- @class CCWSUI
--- @field backend string The backend WebSocket URL to connect to.
--- @field textures table<string, string> A table of texture IDs to base64 encoded images.
--- @field private state table<string, table> The current state of the UI.
--- @field package handlers table<string, table<string, CCWSUI.Handler>>
--- @field private ws table The WebSocket connection to the backend.
--- @field render fun(self: CCWSUI, ctx: CCWSUI.Context): CCWSUI.Component A function that renders the UI.
--- @field ready fun(self: CCWSUI, url: string)
--- @field private rerendering boolean Whether the UI is currently being rerendered.
--- @field error string|nil
local CCWSUI = {}
CCWSUI.__index = CCWSUI

--- Create a new CCWSUI instance. The value of wantedslug may not be the slug
--- that is actually used. Listen to CCWSUI.ready to get the actual URL for the
--- user to visit, which contains the actual slug.
--- @param wantedslug string|nil The slug to use for the UI.
--- @param backend string|nil The backend WebSocket URL to connect to.
--- @return CCWSUI
function CCWSUI.new(wantedslug, backend)
	expect(1, backend, "string", "nil")
	backend = backend or "ws://localhost:8080"
	local ws, err = http.websocket(backend .. "/host")
	if not ws then error(err, 2) end
	if wantedslug then
		ws.send(textutils.serializeJSON({
			t = 2, d = { slug = wantedslug }
		}))
	end
	return setmetatable({
		backend = backend,
		ws = ws,
		textures = {},
		state = {},
		handlers = {},
	}, CCWSUI)
end

--- Run runs the CCWSUI main loop. Place this in a parallel.waitForAny call.
--- It only returns if something goes wrong. Check CCWSUI.error after the
--- parallel call to see what went wrong, if anything.
function CCWSUI:run()
	if not self.render then error("CCWSUI:render is not set!", 2) end
	if not self.ready then error("CCWSUI:ready is not set!", 2) end
	self.ws.send(textutils.serializeJSON({ t = 3 })) -- freeze
	while true do
		local jmsg, err = self.ws.receive()
		if err then
			self.error = err
			return
		end
		if jmsg == nil then
			self.error = "Connection closed"
			return
		end
		local msg = textutils.unserializeJSON(jmsg)
		if msg.t == 4 then -- Ready
			local secure, domain = self.backend:match("^ws(s?)://([^/]+)")
			local url = "http" .. secure .. "://" .. domain .. msg.d.url
			self:ready(url)
		elseif msg.t == 5 then -- Hello
			local state = {}
			self.state[msg.d.client] = setmetatable({}, {
				__index = state,
				__newindex = function(t, k, v)
					if self.rerendering then
						error("CCWSUI:render() may not modify state (deadlock).", 2)
					end
					state[k] = v
					self:forceRender(msg.d.client)
				end
			})
			self:forceRender(msg.d.client)
		elseif msg.t == 6 then -- Leave
			self.state[msg.d.client] = nil
			self.handlers[msg.d.client] = nil
		elseif msg.t == 7 then -- Event
			os.queueEvent("ccwsui:event", msg.d.client, msg.d.event, msg.d.data)
			if self.handlers[msg.d.client] and self.handlers[msg.d.client][msg.d.event] then
				self.handlers[msg.d.client][msg.d.event](msg.d.data)
			end
		else
			self.error = "Got unknown message type: " .. msg.t
			return
		end
	end
end

--- Gets the state for a specific client.
--- @param client string The client ID to get the state for.
--- @return table|nil
function CCWSUI:getState(client) return self.state[client] end

--- @class CCWSUI.Context
--- @field client string The client ID.
--- @field state table The state for the client.
--- @field private inst CCWSUI
local Context = {}

local function contextFrom(inst, client)
	return setmetatable({
		client = client,
		state = inst.state[client],
		inst = inst
	}, { __index = Context })
end

--- @param event string The unique event to add a handler for.
--- @param handler CCWSUI.Handler The handler to add.
function Context:addHandler(event, handler)
	local h = self.inst.handlers[self.client]
	if h[event] then error("Handler already exists for event " .. event) end
	h[event] = handler
end

--- Forces a render for a specific client.
--- @param client string The client ID to force render for.
function CCWSUI:forceRender(client)
	expect(1, client, "string")
	self.rerendering = true
	self.handlers[client] = {}
	local root = self:render(contextFrom(self, client))
	self.rerendering = false
	self.ws.send(textutils.serializeJSON({
		t = 1, -- Update
		d = { client = client, root = root }
	}))
end

return CCWSUI
