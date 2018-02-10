# runelite-launcher

## Building

```shell
mvn clean install
```

This will build distribution for all major platforms (linux32, linux64, windows32, windows64, darwin)
and runelite-launcher for your current platform.

### Cross-compiling GUI launcher

For cross-compilation of GUI launcher, various libraries and compiler toolchains are needed, mainly:

* Standard GNU GCC
* MinGW64
* osxcross (or being on osx)

To cross-compile GUI (linux32, linux64, windows32, windows64, darwin), run

```
mvn clean install -Pgui
```

### Cross-compiling TUI launcher

Cross-compiling TUI launcher is a lot simplier, as all we need is Go compiler and nothing else.
To cross-compile TUI (linux32, linux64, windows32, windows64, darwin), run

```
mvn clean install -Ptui
```

## Launching the launcher

```
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT
```
