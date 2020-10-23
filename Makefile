.PHONY: precommit test lint_cmd test_cmd build_cmd release_cmd clean
test: test_cmd plan_terraform

lint_cmd: precommit
	make -C decider lint

test_cmd: precommit
	make -C decider test

build_cmd: precommit
	make -C decider build

release_cmd: precommit
	make -C decider release

clean: precommit
	make -C decider clean

precommit:
ifneq ($(strip $(hooksPath)),.github/hooks)
	@git config --add core.hooksPath .github/hooks
endif
