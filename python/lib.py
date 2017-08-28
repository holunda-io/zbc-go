from sys import platform as _platform
from ctypes import cdll

def lib_extension():
    if _platform == 'darwin':
        return 'dylib'
    elif _platform == 'win32':
        return 'dll'
    else:
        return 'so'

lib = cdll.LoadLibrary('lib/zbc.' + lib_extension())

print "Python say: about to call Go ..."
total = lib.Add(40, 2)
print "Python says: result is " + str(total)

client, err = lib.NewClient("0.0.0.0:51015")
print "Python says: result is " + str(client)