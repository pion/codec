#pragma once

#include <openh264/codec_api.h>

#ifdef __cplusplus
extern "C"
{
#endif
    typedef struct Slice
    {
        unsigned char *data;
        int data_len;
    } Slice;

    typedef struct Frame
    {
        void *y, *u, *v;
        int height;
        int width;
    } Frame;

    typedef struct Encoder
    {
        SEncParamBase params;
        ISVCEncoder *engine;
        unsigned char *buff;
    } Encoder;

    Encoder *enc_new(const SEncParamBase params);
    void enc_free(Encoder *e);
    Slice enc_encode(Encoder *e, Frame f);
#ifdef __cplusplus
}
#endif