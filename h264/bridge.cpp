#include "bridge.hpp"
#include <errno.h>
#include <stdlib.h>
#include <string.h>

Encoder *enc_new(const SEncParamBase params)
{
    int rv;
    ISVCEncoder *engine;

    rv = WelsCreateSVCEncoder(&engine);
    if (rv != 0)
    {
        errno = rv;
        return NULL;
    }

    rv = engine->Initialize(&params);
    if (rv != 0)
    {
        errno = rv;
        return NULL;
    }

    Encoder *encoder = (Encoder *)malloc(sizeof(Encoder));
    encoder->engine = engine;
    encoder->params = params;
    encoder->buff = NULL;
    return encoder;
}

void enc_free(Encoder *e)
{
    int rv = e->engine->Uninitialize();
    if (rv != 0)
    {
        errno = rv;
        return;
    }

    WelsDestroySVCEncoder(e->engine);

    free(e);
}

// There's a good reference from ffmpeg in using the encode_frame
// Reference: https://ffmpeg.org/doxygen/2.6/libopenh264enc_8c_source.html
Slice enc_encode(Encoder *e, Frame f)
{
    int rv;
    SSourcePicture pic;
    SFrameBSInfo info;

    // Make sure that we clear up previous buff
    free(e->buff);

    pic.iPicWidth = f.width;
    pic.iPicHeight = f.height;
    pic.iColorFormat = videoFormatI420;
    // Since we're using 4:2:0 format, we can set the stride for the chromas
    // to be the the width of the frame. That way we can skip even rows.
    // For example, if we have a picture of 400x400, our chroma will be 400x200
    // from Go.
    pic.iStride[0] = pic.iStride[1] = pic.iStride[2] = pic.iPicWidth;
    pic.pData[0] = (unsigned char *)f.y;
    pic.pData[1] = (unsigned char *)f.u;
    pic.pData[2] = (unsigned char *)f.v;

    rv = e->engine->EncodeFrame(&pic, &info);
    if (rv != 0)
    {
        errno = rv;
        return Slice{0};
    }

    int *layer_size = (int *)calloc(sizeof(int), info.iLayerNum);
    int size = 0;
    for (int layer = 0; layer < info.iLayerNum; layer++)
    {
        for (int i = 0; i < info.sLayerInfo[layer].iNalCount; i++)
            layer_size[layer] += info.sLayerInfo[layer].pNalLengthInByte[i];

        size += layer_size[layer];
    }

    e->buff = (unsigned char *)malloc(size);
    size = 0;
    for (int layer = 0; layer < info.iLayerNum; layer++)
    {
        memcpy(e->buff + size, info.sLayerInfo[layer].pBsBuf, layer_size[layer]);
        size += layer_size[layer];
    }
    free(layer_size);

    Slice s = {.data = e->buff, .data_len = size};
    return s;
}