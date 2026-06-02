local CCWSUI = require("ccwsui").new()

function CCWSUI:render(client)
	return {
		kind = "literal",
		props = {
			Pieces = { { Text = "Hello from Lua!", Color = { R = 255, G = 0, B = 0, A = 255 } } }
		},
	}
end

function CCWSUI:ready(url)
	print("Access UI: " .. url)
end

CCWSUI:run()
