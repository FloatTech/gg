kernel void scale(
    read_only image2d_t inputImg,
    sampler_t smp,
    write_only image2d_t outputImg)
{
    uint x = get_global_id(0);
    uint y = get_global_id(1);
    uint outW = get_image_width(outputImg);
    uint outH = get_image_height(outputImg);

    if (x >= outW || y >= outH) {
        return;
    }

    float2 normCoord = (float2)(
        (float)x / (float)outW,
        (float)y / (float)outH
    );

    float4 pixel = read_imagef(inputImg, smp, normCoord);

    write_imagef(outputImg, (int2)(x, y), pixel);
}
