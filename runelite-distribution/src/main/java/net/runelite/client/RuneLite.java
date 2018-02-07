package net.runelite.client;

public class RuneLite
{
    public static void main(String... args)
    {
        System.err.println(
            "You are trying to launch dummy launcher.\n" +
                "This means that you probably do not have correctly updated RuneLite client.\n" +
                "Try to clear ~/.runelite/cache and run the launcher again.\n" +
                "If nothing changed, please file an issue here: https://github.com/runelite/launcher");

        System.exit(1);
    }
}
