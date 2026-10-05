# Hints - 10

<details>
<summary>Hint 1 - direction</summary>

Each fake's `Notify` is guard → guard → append. `LoggingNotifier.Notify` is: log, nil-check, delegate.
`NotifyAll` is a loop that *collects* errors instead of returning on the first one.

</details>

<details>
<summary>Hint 2 - what to look at</summary>

- [`strings.Contains`](https://pkg.go.dev/strings#Contains) for `@`, [`strings.HasPrefix`](https://pkg.go.dev/strings#HasPrefix) for `+`.
- `len(msg) > SMSMaxLen` for the length (bytes are fine here).
- Return `fmt.Errorf("sms to %q: %w", s.To, ErrInvalidAddress)` or the bare sentinel — both pass.
- The embedded field's name is its type name: `l.Notifier`. Delegate with `l.Notifier.Notify(msg)`.
- Append to `l.Log` **before** the nil check and before delegating.
- `NotifyAll`: build a `[]error`, `append` each non-nil error, then `return errors.Join(errs...)`.
  `errors.Join` with no non-nil arguments returns `nil`, so no special case is needed.

</details>

<details>
<summary>Hint 3 - near-pseudocode</summary>

```
(e *EmailNotifier) Notify(msg):
    if !Contains(e.To, "@"): return ErrInvalidAddress
    e.Sent = append(e.Sent, msg); return nil

(s *SMSNotifier) Notify(msg):
    if !HasPrefix(s.To, "+"):  return ErrInvalidAddress
    if len(msg) > SMSMaxLen:   return ErrMessageTooLong
    append; return nil

(l *LoggingNotifier) Notify(msg):
    l.Log = append(l.Log, msg)
    if l.Notifier == nil: return ErrNoNotifier
    return l.Notifier.Notify(msg)

NotifyAll(ns, msg):
    var errs []error
    for each n in ns:
        if err := n.Notify(msg); err != nil: errs = append(errs, err)
    return errors.Join(errs...)
```

</details>
