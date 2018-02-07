# launcher-go

WIP, still not working

## Testing

```shell
mvn clean install
mkdir -p ~/.runelite/cache

# Copy distribution to cache, and run launcher that will extract file, simulates
# downloading distribution and unpacking it for now
cp runelite-distribution/target/runelite-distribution-1.0.0-SNAPSHOT-archive-distribution-linux64.tar.gz ~/.runelite/cache
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT

# Remove .jar from distribution, simulates having outdated client, and then run
# launcher again this will download correct shaded jar from RuneLite repository
# and now also launches it
rm ~/.runelite/cache/linux64/runelite-distribution-1.0.0-SNAPSHOT/runelite-distribution-1.0.0-SNAPSHOT.jar
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT
```

## Launching the launcher

```
./runelite-launcher/target/runelite-launcher-1.0.0-SNAPSHOT
```
