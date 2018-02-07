# launcher-go

WIP, still not working

## Testing

```shell
mvn clean install
mkdir -p ~/.runelite/cache

# Copy distribution to cache, simulates downloading distribution for now
cp runelite-distribution/target/runelite-distribution-1.0.0-SNAPSHOT-archive-distribution-linux64.tar.gz ~/.runelite/cache
```

## Launching the launcher

```
# Now launcher will unpack distribution and download latest shaded client, and
# launch the packr distribution
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT
```
