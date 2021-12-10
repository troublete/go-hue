-- Turn off all lamps
--
-- Usage:
--	hue 03-turn-off-all-lamps.lua
 
local first = nil
for serial, b in pairs(bridges) do 
	if not first then
		first = b
	end
end

for _, l in pairs(first.lights(os.getenv('HUE_USERNAME'))) do
	l.turn_off()
end
