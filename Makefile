PROJECT_ROOT := src/
VERSION = 0.1

.DEFAULT_GOAL := all

# The resulting binaries map to the subproject names
BINARIES = \
	uspin

LIBRARIES = \
	libuspin/boot \
	libuspin/build \
	libuspin/commands \
	libuspin/config \
	libuspin/disk \
	libuspin/pkg \
	libuspin/spec

GO_TESTS = \
	$(addsuffix .test,$(LIBRARIES))

include Makefile.gobuild

# We want to add compliance for all built binaries
_CHECK_COMPLIANCE = $(addsuffix .compliant,$(BINARIES)) $(addsuffix .compliant,$(LIBRARIES))

# Build all binaries as static binary
BINS = $(addsuffix .statbin,$(BINARIES))

# Ensure our own code is compliant..
compliant: $(_CHECK_COMPLIANCE)
install: $(BINS)
	test -d $(DESTDIR)/usr/bin || install -D -d -m 00755 $(DESTDIR)/usr/bin; \
	install -m 00755 builds/* $(DESTDIR)/usr/bin/.

ensure_modules:
	@ ( \
		git submodule init; \
		git submodule update; \
	);

release:
	git archive --format=tar.gz --verbose -o USpin-$(VERSION).tar.gz HEAD --prefix=USpin-$(VERSION)/

all: $(BINS)
