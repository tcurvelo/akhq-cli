# akhq-cli

A command-line client for AKHQ. Main uses: watching offset lag at a finer interval than metrics scraping allows, fetching every message of a topic (for example a dead letter topic, for reports or manual replay), and finding specific messages.

## Language

**Offset lag**:
The number of messages a consumer group has not yet consumed on a topic: per partition, the latest offset minus the group's committed offset.
_Avoid_: Lag (unqualified), offsets, time lag

**Partial offset lag**:
An offset lag total that leaves out partitions where the consumer group has never committed an offset, so the true lag is unknown and may be higher.
_Avoid_: Lag of 0 for uncommitted partitions

**Topic dump**:
Every message of a topic, or of a time range within it, fetched oldest first; typically a dead letter topic, for a report or a manual replay.
_Avoid_: Export, read all
