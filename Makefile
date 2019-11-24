vendor_dir = cvendor
src_dir = $(vendor_dir)/src
lib_dir = $(vendor_dir)/lib
include_dir = $(vendor_dir)/include

vendor: create_out_dirs openh264 x264 

openh264: $(src_dir)/openh264
	cd $^ && $(MAKE)
	mkdir -p $(lib_dir)/$@ $(include_dir)/$@
	cp $^/libopenh264.a $(lib_dir)/$@
	cp $^/codec/api/svc/*.h $(include_dir)/$@

x264: $(src_dir)/x264
	cd $^ && $(MAKE)
	mkdir -p $(lib_dir)/$@ $(include_dir)/$@
	cp $^/libx264.a $(lib_dir)/$@
	cp $^/x264_config.h $(include_dir)/$@
	cp $^/x264.h $(include_dir)/$@

create_out_dirs:
	mkdir -p $(lib_dir) $(include_dir)

dependencies:
	apt install -y nasm