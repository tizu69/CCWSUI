local expect = require "cc.expect"

--- @alias CCWSUI.Handler<T> fun(ev: T)
--- @alias CCWSUI.StateProxy table<string, any>

--- @class CCWSUI
--- @field private wantedslug string|nil
--- @field backend string The backend WebSocket URL to connect to.
--- @field textures table<string, string> A table of texture IDs to base64 encoded images.
--- @field private users table<string, string>
--- @field private state table<string, CCWSUI.StateProxy>
--- @field package handlers table<string, table<string, {handler: CCWSUI.Handler; timeout: number}>>
--- @field private ws table The WebSocket connection to the backend.
--- @field render fun(self: CCWSUI, ctx: CCWSUI.Context): CCWSUI.Component A function that renders the UI.
--- @field hello fun(self: CCWSUI, client: string, user: string) Gets called when a new client connects.
--- @field leave fun(self: CCWSUI, client: string) Gets called when a client disconnects.
--- @field handle fun(self: CCWSUI, client: string, event: string, data: table) Gets called when a client sends an event.
--- @field ready fun(self: CCWSUI, url: string) Gets called when the UI is ready to be used. May be called more than once if the connection drops.
--- @field private rerendering boolean Whether the UI is currently being rerendered.
--- @field private retryDelay number
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
	backend = backend or "wss://ccwsui.tizu.dev"
	local ws, err = http.websocket(backend .. "/host")
	if not ws then error(err, 2) end
	if wantedslug then
		ws.send(textutils.serializeJSON({
			t = 2, d = { slug = wantedslug }
		}))
	end
	return setmetatable({
		backend = backend,
		wantedslug = wantedslug,
		ws = ws,
		textures = {},
		users = {},
		state = {},
		handlers = {},
		retryDelay = 0,
	}, CCWSUI)
end

--- Run runs the CCWSUI main loop. Place this in a parallel.waitForAny call.
--- It only returns if something goes wrong. Check CCWSUI.error after the
--- parallel call to see what went wrong, if anything.
function CCWSUI:run()
	self.retryDelay = 0
	self.ws.send(textutils.serializeJSON({ t = 3 })) -- freeze
	while true do
		local pcallok, jmsg, err = pcall(self.ws.receive)
		if not pcallok then
			self.error = jmsg
			return
		elseif err then
			self.error = err
			return
		elseif jmsg == nil then
			self.error = "Connection closed"
			self:reconnect()
		elseif jmsg ~= "" then -- keepalive
			self.retryDelay = 0
			local msg = textutils.unserializeJSON(jmsg)
			if msg.t == 4 then -- Ready
				local secure, domain = self.backend:match("^ws(s?)://([^/]+)")
				local url = "http" .. secure .. "://" .. domain .. msg.d.url
				self:ready(url)
			elseif msg.t == 5 then -- Hello
				self.users[msg.d.client] = msg.d.user
				self:hello(msg.d.client, msg.d.user)
				os.queueEvent("ccwsui:hello", msg.d.client, msg.d.user)
				self:forceRender(msg.d.client)
			elseif msg.t == 6 then -- Leave
				self:leave(msg.d.client)
				os.queueEvent("ccwsui:leave", msg.d.client)
				for _, state in pairs(self.state) do
					state._listeners[msg.d.client] = nil
				end
				self.users[msg.d.client] = nil
				self.handlers[msg.d.client] = nil
			elseif msg.t == 7 then -- Event
				self:handle(msg.d.client, msg.d.event, msg.d.data)
				os.queueEvent("ccwsui:event", msg.d.client, msg.d.event, msg.d.data)
				if self.handlers[msg.d.client] and self.handlers[msg.d.client][msg.d.event] then
					self.handlers[msg.d.client][msg.d.event].handler(msg.d.data)
				end
			else
				self.error = "Got unknown message type: " .. msg.t
				return
			end
		end
	end
end

--- Reconnects to the backend with exponential backoff (maximum of 30s).
--- On success, re-sends slug and freeze, updates self.ws, and returns.
function CCWSUI:reconnect()
	if self.retryDelay == 0 then self.retryDelay = 1 end
	while true do
		local delay = math.min(self.retryDelay, 30)
		sleep(delay)

		local ws, err = http.websocket(self.backend .. "/host")
		if ws then
			if self.wantedslug then
				ws.send(textutils.serializeJSON({ t = 2, d = { slug = self.wantedslug } }))
			end
			for id, data in pairs(self.textures) do
				ws.send(textutils.serializeJSON({ t = 9, d = { id = id, data = data } }))
			end
			ws.send(textutils.serializeJSON({ t = 3 })) -- freeze
			self.ws = ws
			self.error = nil
			self.retryDelay = 0
			return
		end
		self.error = err
		self.retryDelay = math.min(self.retryDelay * 2, 30)
	end
end

--- Returns the user UUID for the given client UUID.
--- Never ever expose the user ID to the client, nor allow the user ID alone to
--- be the deciding factor for things like authorization. The user ID is not a
--- secret, and can be easily spoofed by a malicious client.
--- @param client string
--- @return string|nil
function CCWSUI:getUser(client) return self.users[client] end

--- Returns whether a state exists for the given key.
--- @param key string
--- @return boolean
function CCWSUI:existsState(key)
	if key == "user" then error("Key 'user' is reserved for internal use.", 2) end
	return self.state[key] ~= nil
end

--- Gets or creates a State proxy for the given key. Clients that are listening
--- to this key will automatically be rerendered when a top-level value is
--- modified.
---
--- ```lua
--- local state = ui:getState("myState")
--- state.foo = "bar" -- rerenders all clients listening to "myState"
--- state.bar.baz = "qux" -- does NOT rerender
--- ```
---
--- The key `user` is reserved for internal use. To match the behavior of ctx.s,
--- use the key `"user:" .. ccwsui:getUser(client)`.
---
--- @param key string
--- @return CCWSUI.StateProxy
function CCWSUI:getState(key)
	if key == "user" then error("Key 'user' is reserved for internal use.", 2) end
	if not self.state[key] then
		local backing = {}
		self.state[key] = setmetatable({
			_backing = backing,
			_listeners = {},
		}, {
			__index = function(t, k)
				if k:sub(1, 1) == "_" then return rawget(t, k) end
				return backing[k]
			end,
			__newindex = function(t, k, v)
				if k:sub(1, 1) == "_" then
					rawset(t, k, v)
					return
				end
				if self.rerendering then
					error("CCWSUI:render() may not modify state (deadlock).", 2)
				end
				backing[k] = v
				os.queueEvent("ccwsui:stateupdate", key, k)
				if t._listeners then
					for client in pairs(t._listeners) do
						self:forceRender(client)
					end
				end
			end,
			__pairs = function() return pairs(backing) end,
		})
		rawset(self.state[key], "_listeners", {})
	end
	return self.state[key]
end

--- @class CCWSUI.Context
--- @field client string The client ID.
--- @field private inst CCWSUI
--- @field user string Unique to a user. May be shared across clients. Never ever expose the user ID to the client, nor allow the user ID alone to be the deciding factor for things like authorization. The user ID is not a secret, and can be easily spoofed by a malicious client.
local Context = {}

local function contextFrom(inst, client)
	return setmetatable({
		client = client,
		inst = inst,
		user = inst:getUser(client),
	}, { __index = Context })
end

--- @param event string The unique event to add a handler for.
--- @param handler CCWSUI.Handler The handler to add.
function Context:addHandler(event, handler)
	local h = self.inst.handlers[self.client]
	if h[event] then error("Handler already exists for event " .. event) end
	h[event] = { handler = handler }
end

--- Returns a State proxy for the given key. As long as this method is called
--- during a render, updates to this state will cause the client to re-render.
---
--- If no key is provided, a default key of `client:<ctx.client>` is used.
--- If the provided key is exactly `user`, `user:<ctx.user>` is used.
---
--- @param key string|nil
--- @return CCWSUI.StateProxy
function Context:s(key)
	key = key or ("client:" .. self.client)
	if key == "user" then key = "user:" .. self.user end
	local s = self.inst:getState(key)
	s._listeners[self.client] = true
	return s
end

--- Forces a render for a specific client.
--- @param client string The client ID to force render for.
function CCWSUI:forceRender(client)
	expect(1, client, "string")

	-- if disconnected, don't try to render
	if self.retryDelay > 0 then return end

	-- reset per-render listener subscriptions; ctx:state() re-populates it
	for _, state in pairs(self.state) do
		state._listeners[client] = nil
	end
	self.rerendering = true
	self.handlers[client] = self.handlers[client] or {}
	for id, handler in pairs(self.handlers[client]) do
		if not handler.timeout then handler.timeout = os.time() + 5 end
		if handler.timeout <= os.time() then self.handlers[client][id] = nil end
	end
	local root = self:render(contextFrom(self, client))
	self.rerendering = false
	self:send(1, { client = client, root = root })
end

--- @param type number The type of message to send.
--- @param data any The data to send.
function CCWSUI:send(type, data)
	pcall(self.ws.send, textutils.serializeJSON({ t = type, d = data }))
end

---@class CCWSUI.Metadata
---@field title string|nil

--- Updates the metadata for a specific client. Unset field are kept unchanged.
--- Use of a zero value (empty string, 0, false, etc) resets the field to its
--- default value on the client side.
--- @param client string The client ID to update metadata for.
--- @return fun(d: CCWSUI.Metadata)
function CCWSUI:updateMetadata(client)
	return function(d)
		self:send(8, { client = client, title = d.title })
	end
end

local alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
local function b64encode(str)
	local lookup = {}
	for i = 1, #alphabet do lookup[i] = string.sub(alphabet, i, i) end

	local len = #str
	local remainder = len % 3
	local out = ""
	for i = 1, len - remainder, 3 do
		local c1, c2, c3 = string.byte(str, i, i + 2)
		out = out ..
			lookup[bit32.rshift(c1, 2) + 1] ..
			lookup[c1 % 4 * 16 + bit32.rshift(c2, 4) + 1] ..
			lookup[c2 % 16 * 4 + bit32.rshift(c3, 6) + 1] ..
			lookup[c3 % 64 + 1]
	end

	if remainder == 2 then
		local c1, c2 = string.byte(str, len - 1, len)
		out = out ..
			lookup[bit32.rshift(c1, 2) + 1] ..
			lookup[c1 % 4 * 16 + bit32.rshift(c2, 4) + 1] ..
			lookup[c2 % 16 * 4 + 1] ..
			"="
	elseif remainder == 1 then
		local c1 = string.byte(str, len)
		out = out .. lookup[bit32.rshift(c1, 2) + 1] .. lookup[c1 % 4 * 16 + 1] .. "=="
	end

	return out
end

--- Uploads a texture for use with components.Texture. The ID must be unique,
--- a duplicate ID will overwrite the previous texture. An ID may also shadow
--- built-in texture IDs, if you want to override any of the built-in textures.
--- Must be called before running `CCWSUI:run()`.
--- @param id string The ID to use for the texture.
--- @param path string A path to the .png file.
function CCWSUI:uploadTexture(id, path)
	local f = fs.open(path, "rb")
	if not f then error("File not found: " .. path, 2) end
	local b64 = b64encode(f.readAll())
	self.textures[id] = b64
	self:send(9, { id = id, data = b64 })
	f.close()
end

-- default hook implementations

function CCWSUI:render(ctx)
	return {
		kind = "align",
		props = { X = 0.5, Y = 0.5 },
		children = {
			{
				kind = "literal",
				props = {
					Wrap = true,
					Pieces = {
						{
							Text = "Hello World! If you're seeing this, congratulations, CCWSUI is working!",
							Color = { R = 255, G = 255, B = 255, A = 255 },
						},
					},
				},
				children = textutils.empty_json_array,
			}
		},
	}
end

function CCWSUI:hello(client, user)
end

function CCWSUI:leave(client)
end

function CCWSUI:handle(client, event, data)
end

function CCWSUI:ready(url)
	print("CCWSUI: " .. url)
end

return CCWSUI
