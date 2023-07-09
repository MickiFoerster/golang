TARGET = hello
TARGET_BPF = $(TARGET).bpf.o

GO_SRC = *.go
BPF_SRC = *.bpf.c

LIBBPF_HEADERS = $$PWD/libbpf/src/root/usr/include
LIBBPF_OBJ = $$PWD/libbpf/src/root/usr/lib64/libbpf.a
GO_ENV = CC=clang CGO_CFLAGS="-I $(LIBBPF_HEADERS)" CGO_LDFLAGS="$(LIBBPF_OBJ)"

.PHONY: all
all: $(TARGET) $(TARGET_BPF) libbpf

$(TARGET): $(GO_SRC)
	$(GO_ENV) go build -o $@

$(TARGET_BPF): $(BPF_SRC)
	clang \
		-I /usr/include/x86_64-linux-gnu \
		-I $(LIBBPF_HEADERS) \
		-O2 \
		-c  \
		-target bpf \
		-o $@ \
		$^


.PHONY: clean
clean:
	rm -rf libbpf
	go clean

libbpf:
	git clone https://github.com/libbpf/libbpf.git
	cd libbpf/src && \
	mkdir build root && \
	BUILD_STATIC_ONLY=y OBJDIR=build DESTDIR=root make install
