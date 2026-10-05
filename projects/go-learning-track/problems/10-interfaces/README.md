# 10 - Notifiers

**New concept:** interfaces — implicit satisfaction, small interfaces, embedding instead of inheritance
**Builds on:** 05, 06, 08

## Task

`Notifier` is a one-method interface. Two fakes record messages instead of sending them; a decorator
logs; a fan-out function sends through all of them.

| Member | Behaviour |
|---|---|
| `(*EmailNotifier).Notify(msg)` | `To` must contain `@`, else `ErrInvalidAddress`. On success append `msg` to `Sent` |
| `(*SMSNotifier).Notify(msg)` | `To` must start with `+` → else `ErrInvalidAddress`. `len(msg) > SMSMaxLen` → `ErrMessageTooLong`. Check the address first. On success append to `Sent` |
| `(*LoggingNotifier).Notify(msg)` | Append `msg` to `Log` **always**, then delegate to the embedded `Notifier` and return its error. No embedded notifier → `ErrNoNotifier`, no panic |
| `NotifyAll(notifiers, msg)` | Call every notifier, even after failures. Return every failure combined (`errors.Join`), or `nil` |

Failures never append to `Sent`. Wrapping errors with context is fine as long as `errors.Is` still works.

## Acceptance criteria

| # | Criterion | Checked by |
|---|---|---|
| 1 | All three types satisfy `Notifier` | compile-time `var _ Notifier = ...` in the test file |
| 2 | Email: records on success, rejects an address without `@` | `TestEmailNotify`, `TestEmailInvalidAddress` |
| 3 | SMS: records, rejects a missing `+`, enforces the inclusive length limit | `TestSMSNotify`, `TestSMSInvalidAddress`, `TestSMSTooLong` |
| 4 | Logging: logs then delegates; logs even on failure; passes the inner error through | `TestLoggingNotifierLogsAndDelegates`, `TestLoggingNotifierLogsFailures` |
| 5 | Logging with nothing embedded returns `ErrNoNotifier` | `TestLoggingNotifierNoInner` |
| 6 | `NotifyAll` reaches everyone, keeps going after failures, joins all errors | `TestNotifyAllSendsToEveryone`, `TestNotifyAllContinuesAfterFailure`, `TestNotifyAllJoinsErrors` |
| 7 | `NotifyAll(nil, ...)` is `nil` | `TestNotifyAllEmpty` |
| 8 | A type declared in the test file works with `NotifyAll` without mentioning `Notifier` | `TestImplicitSatisfaction` |

## Run

```bash
go test ./problems/10-interfaces/...
go run ./problems/10-interfaces
```

## Stretch (optional)

Once everything passes, comment out `(*LoggingNotifier).Notify` and run the tests again. Everything
still compiles — but the logging tests fail and `TestLoggingNotifierNoInner` panics. Explain in one
line which `Notify` is being called now, why the log is empty, and where the panic comes from. Then
restore the method.
