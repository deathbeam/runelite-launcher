# runelite-launcher

[![](https://jitpack.io/v/deathbeam/runelite-launcher.svg)](https://jitpack.io/#deathbeam/runelite-launcher)

![clean](https://i.imgur.com/HchTzSG.png)

## Running

Download runelite-launcher for your platform from the
[releases page](https://github.com/deathbeam/runelite-launcher/releases).

## Building

```shell
mvn clean install
```

This will build distribution for all major platforms (linux32, linux64, windows32, windows64, darwin)
and runelite-launcher for your current platform.

#### Launching the built launcher

```
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT
```

#### Cross-compiling launcher

For cross-compilation of GUI launcher, you will need [Docker](https://www.docker.com/).  
Now, you need to pull docker image maven will use for cross-compilation:

```
docker pull karalabe/xgo-latest
```

Now, you also need to install [NSIS](http://nsis.sourceforge.net/Main_Page) for
compiler to be able to create Windows Installer. You also need to amke sure that
after installation of NSIS, the `makensis` executable is on your `PATH`.

Now, to actually cross-compile to all platforms (linux32, linux64, windows32, windows64, darwin), run

```
mvn clean install -Pcross
```
