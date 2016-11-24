PROJECT_ROOT := src/

.DEFAULT_GOAL := all

# The resulting binaries map to the subproject names
BINARIES = \
	solspin

LIBRARIES = \
	libspin \
	libspin/config \
	libspin/spec

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

all: $(BINS)
