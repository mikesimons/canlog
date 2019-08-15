# Canlog - canonical log collector

Canlog is a thread safe canonical log collector. Canonical logs are basically structured summaries of events / requests.

## How to use

Upon receiving a request or event, grab a refrence with `canlog.Ref()`. This is a unique string that identifies the event / request and should be propagated through your callstack with context.

When you have information to add to canlog, simply set it with the ref using `canlog.Push(ref, name, value)`.

At the exit point of your request you can generate a standard `map[string]interface{}` from a ref with `canlog.Pop(ref)`.
As the name suggests, this will also remove all references to data for `ref` from canlog.

At this point you are free to add any additional data or emit through your normal logger.

## How it works
Canlog creates a buffered channel for each ref with a reader routine setting values in to a normal golang map.
This mechanism serializes writes to the map thus being thread safe without locks.

When calling `canlog.Pop` the reader routine is instructed to stop. A `sync.WaitGroup` is used to coordinate the reader stopping with the `canlog.Pop` returning.

## License

MIT. See [LICENSE](LICENSE)
