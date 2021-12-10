-- List all the available lamps
--
-- Usage:
--	hue 02-list-all-lamps.lua
 
local first = nil
for serial, b in pairs(bridges) do 
	if not first then
		first = b
	end
end

for id, l in pairs(first.lights(os.getenv('HUE_USERNAME'))) do
	print(id, l.type, l.name) -- prints light id (necessary to set state), type of lamp and name
end
