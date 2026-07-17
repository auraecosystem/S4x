# ==============================================================================
# Modular Assembly Targets: Integrated Sub-Compiler Routine
# ==============================================================================

ASSEMBLY_SRC = $(wildcard native/*.yasm)
ASSEMBLY_OBJ = $(patsubst native/%.yasm, build/obj/%.o, $(ASSEMBLY_SRC))

# Assembler compiler binary selector
YASM_COMPILER = yasm
YASM_FLAGS_LINUX = -f elf64 -D__LINUX__
YASM_FLAGS_WIN   = -f win64 -D__WINDOWS__

.PHONY: compile-asm-layers

## compile-asm-layers: Translate naked vector assembly streams to object code blocks
compile-asm-layers: $(ASSEMBLY_OBJ)
	@echo "--> High-throughput vector architecture mapping completed."

build/obj/%.o: native/%.yasm
	@mkdir -p $(dir $@)
	@echo "--> Processing hardware assembly matrix instructions: $<"
	$(YASM_COMPILER) $(YASM_FLAGS_LINUX) -o $@ $<
