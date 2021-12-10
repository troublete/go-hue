-- Set all lamps to red
--
-- Usage:
--	hue 05-red.lua
 
local first = nil
for serial, b in pairs(bridges) do 
	if not first then
		first = b
	end
end

for _, l in pairs(first.lights(os.getenv('HUE_USERNAME'))) do
	-- param order: hue, brightness, saturation, color temperature, tansition time (in multiple of 100ms);
	-- to ignore param, set to -1; if a lamp doesn't support a property it will just ignore it
	l.set(65500, 255, 255, -1, 10)
end
