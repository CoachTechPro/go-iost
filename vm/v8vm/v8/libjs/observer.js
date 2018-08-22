"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * Provides simple way to "proxify" nested objects and validate the changes.
 */
module.exports = (function () {
    var _create = function(target, path) {
        var proxies = {};
        var getPath = function getPath(path, prop) {
            if (path.length !== 0)
                return path + "." + prop;
            else
                return prop;
        };

        var handler = {
            get: function(target, property, receiver) {
                // var aa = {
                //     target: target,
                //     prop: property,
                //     path: getPath(path, property),
                //     type: typeof target[property]
                // }
                // _native_log('get: ' + JSON.stringify(aa));

                if (target[property] instanceof BigNumber || target[property] instanceof Int64 || typeof target[property] === 'string') {
                    var objectStorage = IOSTContractStorage.get(property, target[property]);
                    return objectStorage;
                }

                var value = Reflect.get(target, property, receiver);
                if (typeof target[property] === 'object' && target[property] !== null) {
                    var objectStorage = IOSTContractStorage.get(property);
                    var proxy = _create(objectStorage, getPath(path, property));
                    proxies[property] = proxy;
                    return proxy;
                }

                if ((property === "set" || property === "get") && typeof value === "function") {
                    const origSet = value;
                    var args = [];
                    for (var _i = 0; _i < arguments.length; _i++) {
                        args[_i] = arguments[_i];
                    }
                    value = function(...args) {
                        console.log(...args);

                        return origSet.apply(target, arguments);
                    };
                }
                return value;
            },
            set: function(target, prop, value, receiver) {
                // var aa = {
                //     target: target,
                //     prop: prop,
                //     path: getPath(path, prop),
                //     value: value,
                //     type: typeof target[prop]
                // }
                // _native_log('set: ' + JSON.stringify(aa));
                target[prop] = value;

                var totalPath = getPath(path, prop);
                var dotIndex = totalPath.indexOf('.');

                if (dotIndex !== -1) {
                    IOSTContractStorage.put(totalPath.substr(0, totalPath.indexOf('.')), target);
                } else {
                    IOSTContractStorage.put(prop, value);
                }
                return true;
            },
        };

        return new Proxy(target, handler);
    };

    return {
        create: function(target) {
            return _create(target, '');
        }
    }
})();