# graalvm-java-cgo-test

Build
---

GraalVM 20.1.0 Java11 and `native-image` are required.

    $ make target/call_from_go

To run

    $ export LD_LIBRARY_PATH=target/native
    $ ./target/call_from_go

