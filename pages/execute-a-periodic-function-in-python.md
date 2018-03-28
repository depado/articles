title: Execute a function periodically using a wrapper in Python
description: Not much to say here, just a snippet to share.
slug: execute-a-periodic-function-in-python
date: 2015-09-04 10:08:00
tags:
    - python
    - snippet

```python
def periodic(interval, times = -1):
    def outer_wrap(function):
        def wrap(*args, **kwargs):
            stop = threading.Event()
            def inner_wrap():
                i = 0
                while i != times and not stop.isSet():
                    stop.wait(interval)
                    function(*args, **kwargs)
                    i += 1

            t = threading.Timer(0, inner_wrap)
            t.daemon = True
            t.start()
            return stop
        return wrap
    return outer_wrap
```
