# runelite-launcher

![clean](https://i.imgur.com/HchTzSG.png)

## Building

```shell
mvn clean install
```

This will build distribution for all major platforms (linux32, linux64, windows32, windows64, darwin)
and runelite-launcher for your current platform.

### Cross-compiling launcher

For cross-compilation of GUI launcher, you will need [Docker](https://www.docker.com/).  
Now, you need to pull docker image maven will use for cross-compilation:

```
docker pull karalabe/xgo-latest
```

Now, to actually cross-compile to all platforms (linux32, linux64, windows32, windows64, darwin), run

```
mvn clean install -Pcross
```

## Launching the launcher

```
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT
```
