from cffi import FFI

ffibuilder = FFI()

ffibuilder.set_source("pyzbc", "", extra_objects=["zbc.so"])

if __name__ == "__main__":
    ffibuilder.compile(verbose=True)
