.PHONY: full

full:
	$(MAKE) -C metrics-audit full
	$(MAKE) -C metrics-traffic full
