XMS = 2g
XMX = 7g

TARGET_CLASS=src/libjavacgo/impl/LibJavaCgo.class

.PHONY: all
all: target/native/libjavacgo.so

.PHONY: clean
clean:
	rm -rf $(TARGET_CLASS)

$(TARGET_CLASS): src/libjavacgo/impl/LibJavaCgo.java
	javac \
	    -cp $(GRAALVM_HOME)/lib/svm/builder/svm.jar \
	    src/libjavacgo/impl/LibJavaCgo.java

target/native/libjavacgo.so: \
	$(TARGET_CLASS)
	mkdir -p target/native
	native-image \
	-cp src \
	-H:Name=libjavacgo \
	--shared \
	-H:+ReportExceptionStackTraces \
	-H:Log=registerResource: \
	-H:+RemoveSaturatedTypeFlows \
	-H:+PrintClassInitialization \
	-H:+TraceClassInitialization \
	--verbose \
	--no-fallback \
	--no-server \
	--initialize-at-build-time \
	$(OPTS) \
	-J-Xms$(XMS) \
	-J-Xmx$(XMX)
	mv graal_isolate_dynamic.h graal_isolate.h libjavacgo.h libjavacgo.so libjavacgo_dynamic.h target/native

target/call_from_go: \
	target/native/libjavacgo.so
	go build -o target/call_from_go --ldflags "-s -w -linkmode 'external'" ./example/main.go
