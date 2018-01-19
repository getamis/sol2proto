// Copyright 2017 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

const (
	cmdCAS = "cas"
	cmdCAD = "cad"
)

func luaScript() string {
	return luaScriptStr
}

const luaScriptStr = `
-- This lua script implements CAS based commands using lua and redis commands.

if #KEYS > 0 then error('No Keys should be provided') end
if #ARGV <= 0 then error('ARGV should be provided') end

local command_name = assert(table.remove(ARGV, 1), 'Must provide a command')

local decode = function(val)
    return cjson.decode(val)
end

local encode = function(val)
    return cjson.encode(val)
end

local exists = function(key)
    return redis.call('exists', key) == 1
end

local get = function(key)
    return redis.call('get', key)
end

local setex = function(key, val, ex)
    if ex == "0" then
        return redis.call('set', key, val)
    end
    return redis.call('set', key, val, 'ex', ex)
end

local del = function(key)
    return redis.call('del', key)
end

-- cas is compare-and-swap function which compare the old value's signature
-- if they are the same, then swap with new val
-- noted that $old and $new are json formatted strings
-- and key is keyed with 'revision'
local revision = "Revision"
local cas = function(key, old, new, ttl)
    if not exists(key) then
        error("redis: key is not found")
    end
    local decodedOrig = decode(get(key))
    local decodedOld = decode(old)
    if decodedOrig[revision] == decodedOld[revision] then
        setex(key, new, ttl)
        return "OK"
    else
        error("redis: value has been changed")
    end
end

-- cad is compare-and-del function which compare the old value's signature
-- if they are the same, then the key will be deleted
-- noted that $old is a json formatted string
-- and key is keyed with 'revision'
local cad = function(key, old)
    if not exists(key) then
        error("redis: key is not found")
    end
    local decodedOrig = decode(get(key))
    local decodedOld = decode(old)
    if decodedOrig[revision] == decodedOld[revision] then
        del(key)
        return "OK"
    else
        error("redis: value has been changed")
    end
end

-- Launcher exposes interfaces which be called by passing the arguments.
local Launcher = {
    cas = cas,
    cad = cad
}

local command = assert(Launcher[command_name], 'Unknown command ' .. command_name)
return command(unpack(ARGV))
`
