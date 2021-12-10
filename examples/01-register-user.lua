-- You must create a developer account first to be able to communicate with
-- the bridge
--
-- Usage:
--	hue 01-register-user.lua

local first = nil
for serial, b in pairs(bridges) do 
	if not first then
		first = b
	end

	print(serial, b.base_url) -- print bridge serial and base_url
end

local username = first.register_user('funky-app-name')
print()
print(string.format("username= %s", username)) -- returns username (to be used for all api requests)
print(string.format("export HUE_USERNAME=%s", username)) -- print command to set as env var; run in shell
