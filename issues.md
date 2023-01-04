# Issues
This file will list out the issues we went through developing this plugin.

## Issue 1: Sending Register Message to StreamDeck Does Not Seem to Work
### Symptom
It seems as though sending the `Register` message to the StreamDeck causes our plugin to exit.

Here is what we see in the StreamDeck logs (`~/Library/Logs/ElgatoStreamDeck/StreamDeck0.log`):
```
09:18:13.8118 void ESDCustomPlugin::onNativeProcessFinished(): The plugin 'Golang Do Nothing Plugin' exited normally with code 2
09:18:13.8120 void ESDCustomPlugin::restartNativeProcess(): Restarting plugin 'Golang Do Nothing Plugin' in 0 seconds(s)
09:18:13.8332 void ESDCustomPlugin::onNativeProcessFinished(): The plugin 'Golang Do Nothing Plugin' exited normally with code 2
09:18:13.8333 void ESDCustomPlugin::restartNativeProcess(): Restarting plugin 'Golang Do Nothing Plugin' in 60 seconds(s)
```

Here is what we see in our plugin logs ($TMPDIR/streamdeck-godonothing-log-*):
```
2023/01/04 09:18:13 Starting Golang DoNothing StreamDeck Plugin
2023/01/04 09:18:13 Command line arguments: -port, 28196, -pluginUUID, 480E2E60AD34E4058637650D90E93136, -registerEvent, registerPlugin, -info, {"application":{"font":".AppleSystemUIFont","language":"en","platform":"mac","platformVersion":"13.1.0","version":"6.0.2.17735"},"colors":{"buttonPressedBackgroundColor":"#303030FF","buttonPressedBorderColor":"#646464FF","buttonPressedTextColor":"#969696FF","disabledColor":"#007AFF7F","highlightColor":"#007AFFFF","mouseDownColor":"#2EA8FFFF"},"devicePixelRatio":1,"devices":[{"id":"719F009FF3A55C4D50DC959B1CECC947","name":"Stream Deck Mini","size":{"columns":3,"rows":2},"type":1}],"plugin":{"uuid":"com.mnelis.godonothing","version":"0.1"}},
2023/01/04 09:18:13 Connected to StreamDeck server
2023/01/04 09:18:13 Handlers setup
2023/01/04 09:18:13 Registering
```

### Solution
No solution yet.

