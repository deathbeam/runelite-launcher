# runelite-launcher

## Building

```shell
mvn clean install
```

This will build distribution for all major platforms (linux32, linux64, windows32, windows64, darwin)
and runelite-launcher for your current platform.

### Cross-compiling launcher

For cross-compilation of launcher, various libraries and compiler toolchains are needed, mainly:

* Gtk3 (both 32 and 64 bit/multilib)
* Standard GNU GCC
* MinGW64

To cross-compile launcher (linux32, linux64, windows32, windows64, darwin), navigate to
`runelite-launcher` and run

```
mvn clean install -Pall
```

## Launching the launcher

```
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT
```
