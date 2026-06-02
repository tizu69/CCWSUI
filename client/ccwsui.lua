local expect = require "cc.expect"

--- @class CCWSUI
--- @field backend string The backend WebSocket URL to connect to.
--- @field textures table<string, string> A table of texture IDs to base64 encoded images.
--- @field private ws table The WebSocket connection to the backend.
--- @field render fun(self: CCWSUI, client: string): CCWSUI.Component A function that renders the UI.
--- @field ready fun(self: CCWSUI, url: string)
--- @field error string|nil
local CCWSUI = {}
CCWSUI.__index = CCWSUI

--- @param backend string|nil The backend WebSocket URL to connect to.
--- @return CCWSUI
function CCWSUI.new(backend)
	expect(1, backend, "string", "nil")
	backend = backend or "ws://localhost:8080"
	local ws, err = http.websocket(backend .. "/host")
	if not ws then error(err, 2) end
	return setmetatable({
		backend = backend,
		ws = ws,
		textures = {},
	}, CCWSUI)
end

--- Run runs the CCWSUI main loop. Place this in a parallel.waitForAny call.
--- It only returns if something goes wrong. Check CCWSUI.error after the
--- parallel call to see what went wrong, if anything.
function CCWSUI:run()
	if not self.render then error("CCWSUI:render is not set!", 2) end
	if not self.ready then error("CCWSUI:ready is not set!", 2) end
	while true do
		local jmsg, err = self.ws.receive()
		if err then
			self.error = err
			return
		end
		local msg = textutils.unserializeJSON(jmsg)
		if msg.t == 1 then -- Ready
			self:ready(msg.d.url)
		elseif msg.t == 2 then -- Hello
			self:forceRender(msg.d.client)
		else
			self.error = "Got unknown message type: " .. msg.t
			return
		end
	end
end

--- @param client string The client ID to force render for.
function CCWSUI:forceRender(client)
	self.ws.send(textutils.serializeJSON({
		t = 1, -- Update
		d = {
			client = client,
			root = self:render(client),
		}
	}))
end

return CCWSUI
