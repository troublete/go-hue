-- Run through all the colors
--
-- Usage:
--	hue -interval=500 06-loop.lua
 
local first = nil

if not first then
	for serial, b in pairs(bridges) do 
		if not first then
			first = b
		end
	end
end

for _, l in pairs(first.lights(os.getenv('HUE_USERNAME'))) do
	l.set((loop*2500)%65535, 255, 255, -1, 1)
end
